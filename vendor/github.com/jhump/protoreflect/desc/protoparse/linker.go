package protoparse

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/internal"
	"github.com/jhump/protoreflect/dynamic"
)

type linker struct {
	files          map[string]*parseResult
	descriptorPool map[*dpb.FileDescriptorProto]map[string]proto.Message
	extensions     map[string]map[int32]string
}

func newLinker(files map[string]*parseResult) *linker {
	return &linker{files: files}
}

func (l *linker) linkFiles() (map[string]*desc.FileDescriptor, error) {
	// First, we put all symbols into a single pool, which lets us ensure there
	// are no duplicate symbols and will also let us resolve and revise all type
	// references in next step.
	if err := l.createDescriptorPool(); err != nil {
		return nil, err
	}

	// After we've populated the pool, we can now try to resolve all type
	// references. All references must be checked for correct type, any fields
	// with enum types must be corrected (since we parse them as if they are
	// message references since we don't actually know message or enum until
	// link time), and references will be re-written to be fully-qualified
	// references (e.g. start with a dot ".").
	if err := l.resolveReferences(); err != nil {
		return nil, err
	}

	// Now we've validated the descriptors, so we can link them into rich
	// descriptors. This is a little redundant since that step does similar
	// checking of symbols. But, without breaking encapsulation (e.g. exporting
	// a lot of fields from desc package that are currently unexported) or
	// merging this into the same package, we can't really prevent it.
	linked, err := l.createdLinkedDescriptors()
	if err != nil {
		return nil, err
	}

	// Now that we have linked descriptors, we can interpret any uninterpreted
	// options that remain.
	for _, r := range l.files {
		fd := linked[r.fd.GetName()]
		if err := l.interpretFileOptions(r, fd); err != nil {
			return nil, err
		}
	}

	return linked, nil
}

func (l *linker) createDescriptorPool() error {
	l.descriptorPool = map[*dpb.FileDescriptorProto]map[string]proto.Message{}
	for _, r := range l.files {
		fd := r.fd
		pool := map[string]proto.Message{}
		l.descriptorPool[fd] = pool
		prefix := fd.GetPackage()
		if prefix != "" {
			prefix += "."
		}
		for _, md := range fd.MessageType {
			if err := addMessageToPool(r, pool, prefix, md); err != nil {
				return err
			}
		}
		for _, fld := range fd.Extension {
			if err := addFieldToPool(r, pool, prefix, fld); err != nil {
				return err
			}
		}
		for _, ed := range fd.EnumType {
			if err := addEnumToPool(r, pool, prefix, ed); err != nil {
				return err
			}
		}
		for _, sd := range fd.Service {
			if err := addServiceToPool(r, pool, prefix, sd); err != nil {
				return err
			}
		}
	}
	// try putting everything into a single pool, to ensure there are no duplicates
	// across files (e.g. same symbol, but declared in two different files)
	type entry struct {
		file string
		msg  proto.Message
	}
	pool := map[string]entry{}
	for f, p := range l.descriptorPool {
		for k, v := range p {
			if e, ok := pool[k]; ok {
				desc1 := e.msg
				file1 := e.file
				desc2 := v
				file2 := f.GetName()
				if file2 < file1 {
					file1, file2 = file2, file1
					desc1, desc2 = desc2, desc1
				}
				node := l.files[file2].nodes[desc2]
				return ErrorWithSourcePos{Pos: node.start(), Underlying: fmt.Errorf("duplicate symbol %s: already defined as %s in %q", k, descriptorType(desc1), file1)}
			}
			pool[k] = entry{file: f.GetName(), msg: v}
		}
	}

	return nil
}

func addMessageToPool(r *parseResult, pool map[string]proto.Message, prefix string, md *dpb.DescriptorProto) error {
	fqn := prefix + md.GetName()
	if err := addToPool(r, pool, fqn, md); err != nil {
		return err
	}
	prefix = fqn + "."
	for _, fld := range md.Field {
		if err := addFieldToPool(r, pool, prefix, fld); err != nil {
			return err
		}
	}
	for _, fld := range md.Extension {
		if err := addFieldToPool(r, pool, prefix, fld); err != nil {
			return err
		}
	}
	for _, nmd := range md.NestedType {
		if err := addMessageToPool(r, pool, prefix, nmd); err != nil {
			return err
		}
	}
	for _, ed := range md.EnumType {
		if err := addEnumToPool(r, pool, prefix, ed); err != nil {
			return err
		}
	}
	return nil
}

func addFieldToPool(r *parseResult, pool map[string]proto.Message, prefix string, fld *dpb.FieldDescriptorProto) error {
	fqn := prefix + fld.GetName()
	return addToPool(r, pool, fqn, fld)
}

func addEnumToPool(r *parseResult, pool map[string]proto.Message, prefix string, ed *dpb.EnumDescriptorProto) error {
	fqn := prefix + ed.GetName()
	if err := addToPool(r, pool, fqn, ed); err != nil {
		return err
	}
	for _, evd := range ed.Value {
		vfqn := fqn + "." + evd.GetName()
		if err := addToPool(r, pool, vfqn, evd); err != nil {
			return err
		}
	}
	return nil
}

func addServiceToPool(r *parseResult, pool map[string]proto.Message, prefix string, sd *dpb.ServiceDescriptorProto) error {
	fqn := prefix + sd.GetName()
	if err := addToPool(r, pool, fqn, sd); err != nil {
		return err
	}
	for _, mtd := range sd.Method {
		mfqn := fqn + "." + mtd.GetName()
		if err := addToPool(r, pool, mfqn, mtd); err != nil {
			return err
		}
	}
	return nil
}

func addToPool(r *parseResult, pool map[string]proto.Message, fqn string, dsc proto.Message) error {
	if d, ok := pool[fqn]; ok {
		node := r.nodes[dsc]
		return ErrorWithSourcePos{Pos: node.start(), Underlying: fmt.Errorf("duplicate symbol %s: already defined as %s", fqn, descriptorType(d))}
	}
	pool[fqn] = dsc
	return nil
}

func descriptorType(m proto.Message) string {
	switch m := m.(type) {
	case *dpb.DescriptorProto:
		return "message"
	case *dpb.DescriptorProto_ExtensionRange:
		return "extension range"
	case *dpb.FieldDescriptorProto:
		if m.GetExtendee() == "" {
			return "field"
		} else {
			return "extension"
		}
	case *dpb.EnumDescriptorProto:
		return "enum"
	case *dpb.EnumValueDescriptorProto:
		return "enum value"
	case *dpb.ServiceDescriptorProto:
		return "service"
	case *dpb.MethodDescriptorProto:
		return "method"
	case *dpb.FileDescriptorProto:
		return "file"
	default:
		// shouldn't be possible
		return fmt.Sprintf("%T", m)
	}
}

func (l *linker) resolveReferences() error {
	l.extensions = map[string]map[int32]string{}
	for _, r := range l.files {
		fd := r.fd
		prefix := fd.GetPackage()
		scopes := []scope{fileScope(fd, l)}
		if prefix != "" {
			prefix += "."
		}
		if fd.Options != nil {
			if err := l.resolveOptions(r, fd, "file", fd.GetName(), proto.MessageName(fd.Options), fd.Options.UninterpretedOption, scopes); err != nil {
				return err
			}
		}
		for _, md := range fd.MessageType {
			if err := l.resolveMessageTypes(r, fd, prefix, md, scopes); err != nil {
				return err
			}
		}
		for _, fld := range fd.Extension {
			if err := l.resolveFieldTypes(r, fd, prefix, fld, scopes); err != nil {
				return err
			}
		}
		for _, ed := range fd.EnumType {
			if err := l.resolveEnumTypes(r, fd, prefix, ed, scopes); err != nil {
				return err
			}
		}
		for _, sd := range fd.Service {
			if err := l.resolveServiceTypes(r, fd, prefix, sd, scopes); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *linker) resolveEnumTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, ed *dpb.EnumDescriptorProto, scopes []scope) error {
	enumFqn := prefix + ed.GetName()
	if ed.Options != nil {
		if err := l.resolveOptions(r, fd, "enum", enumFqn, proto.MessageName(ed.Options), ed.Options.UninterpretedOption, scopes); err != nil {
			return err
		}
	}
	for _, evd := range ed.Value {
		if evd.Options != nil {
			evFqn := enumFqn + "." + evd.GetName()
			if err := l.resolveOptions(r, fd, "enum value", evFqn, proto.MessageName(evd.Options), evd.Options.UninterpretedOption, scopes); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *linker) resolveMessageTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, md *dpb.DescriptorProto, scopes []scope) error {
	fqn := prefix + md.GetName()
	scope := messageScope(fqn, l.descriptorPool[fd])
	scopes = append(scopes, scope)
	prefix = fqn + "."

	if md.Options != nil {
		if err := l.resolveOptions(r, fd, "message", fqn, proto.MessageName(md.Options), md.Options.UninterpretedOption, scopes); err != nil {
			return err
		}
	}

	for _, nmd := range md.NestedType {
		if err := l.resolveMessageTypes(r, fd, prefix, nmd, scopes); err != nil {
			return err
		}
	}
	for _, ned := range md.EnumType {
		if err := l.resolveEnumTypes(r, fd, prefix, ned, scopes); err != nil {
			return err
		}
	}
	for _, fld := range md.Field {
		if err := l.resolveFieldTypes(r, fd, prefix, fld, scopes); err != nil {
			return err
		}
	}
	for _, fld := range md.Extension {
		if err := l.resolveFieldTypes(r, fd, prefix, fld, scopes); err != nil {
			return err
		}
	}
	for _, er := range md.ExtensionRange {
		if er.Options != nil {
			erName := fmt.Sprintf("%s:%d-%d", fqn, er.GetStart(), er.GetEnd()-1)
			if err := l.resolveOptions(r, fd, "extension range", erName, proto.MessageName(er.Options), er.Options.UninterpretedOption, scopes); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *linker) resolveFieldTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, fld *dpb.FieldDescriptorProto, scopes []scope) error {
	thisName := prefix + fld.GetName()
	scope := fmt.Sprintf("field %s", thisName)
	node := r.getFieldNode(fld)
	elemType := "field"
	if fld.GetExtendee() != "" {
		fqn, dsc := l.resolve(fd, fld.GetExtendee(), isMessage, scopes)
		if dsc == nil {
			return ErrorWithSourcePos{Pos: node.fieldExtendee().start(), Underlying: fmt.Errorf("unknown extendee type %s", fld.GetExtendee())}
		}
		extd, ok := dsc.(*dpb.DescriptorProto)
		if !ok {
			otherType := descriptorType(dsc)
			return ErrorWithSourcePos{Pos: node.fieldExtendee().start(), Underlying: fmt.Errorf("extendee is invalid: %s is a %s, not a message", fqn, otherType)}
		}
		fld.Extendee = proto.String("." + fqn)
		// make sure the tag number is in range
		found := false
		tag := fld.GetNumber()
		for _, rng := range extd.ExtensionRange {
			if tag >= rng.GetStart() && tag < rng.GetEnd() {
				found = true
				break
			}
		}
		if !found {
			return ErrorWithSourcePos{Pos: node.fieldTag().start(), Underlying: fmt.Errorf("%s: tag %d is not in valid range for extended type %s", scope, tag, fqn)}
		}
		// make sure tag is not a duplicate
		usedExtTags := l.extensions[fqn]
		if usedExtTags == nil {
			usedExtTags = map[int32]string{}
			l.extensions[fqn] = usedExtTags
		}
		if other := usedExtTags[fld.GetNumber()]; other != "" {
			return ErrorWithSourcePos{Pos: node.fieldTag().start(), Underlying: fmt.Errorf("%s: duplicate extension: %s and %s are both using tag %d", scope, other, thisName, fld.GetNumber())}
		}
		usedExtTags[fld.GetNumber()] = thisName
		elemType = "extension"
	}

	if fld.Options != nil {
		if err := l.resolveOptions(r, fd, elemType, thisName, proto.MessageName(fld.Options), fld.Options.UninterpretedOption, scopes); err != nil {
			return err
		}
	}

	if fld.GetTypeName() == "" {
		// scalar type; no further resolution required
		return nil
	}

	fqn, dsc := l.resolve(fd, fld.GetTypeName(), isType, scopes)
	if dsc == nil {
		return ErrorWithSourcePos{Pos: node.fieldType().start(), Underlying: fmt.Errorf("%s: unknown type %s", scope, fld.GetTypeName())}
	}
	switch dsc := dsc.(type) {
	case *dpb.DescriptorProto:
		fld.TypeName = proto.String("." + fqn)
	case *dpb.EnumDescriptorProto:
		fld.TypeName = proto.String("." + fqn)
		// we tentatively set type to message, but now we know it's actually an enum
		fld.Type = dpb.FieldDescriptorProto_TYPE_ENUM.Enum()
	default:
		otherType := descriptorType(dsc)
		return ErrorWithSourcePos{Pos: node.fieldType().start(), Underlying: fmt.Errorf("%s: invalid type: %s is a %s, not a message or enum", scope, fqn, otherType)}
	}
	return nil
}

func (l *linker) resolveServiceTypes(r *parseResult, fd *dpb.FileDescriptorProto, prefix string, sd *dpb.ServiceDescriptorProto, scopes []scope) error {
	thisName := prefix + sd.GetName()
	if sd.Options != nil {
		if err := l.resolveOptions(r, fd, "service", thisName, proto.MessageName(sd.Options), sd.Options.UninterpretedOption, scopes); err != nil {
			return err
		}
	}

	for _, mtd := range sd.Method {
		if mtd.Options != nil {
			if err := l.resolveOptions(r, fd, "method", thisName+"."+mtd.GetName(), proto.MessageName(mtd.Options), mtd.Options.UninterpretedOption, scopes); err != nil {
				return err
			}
		}
		scope := fmt.Sprintf("method %s.%s", thisName, mtd.GetName())
		node := r.getMethodNode(mtd)
		fqn, dsc := l.resolve(fd, mtd.GetInputType(), isMessage, scopes)
		if dsc == nil {
			return ErrorWithSourcePos{Pos: node.getInputType().start(), Underlying: fmt.Errorf("%s: unknown request type %s", scope, mtd.GetInputType())}
		}
		if _, ok := dsc.(*dpb.DescriptorProto); !ok {
			otherType := descriptorType(dsc)
			return ErrorWithSourcePos{Pos: node.getInputType().start(), Underlying: fmt.Errorf("%s: invalid request type: %s is a %s, not a message", scope, fqn, otherType)}
		}
		mtd.InputType = proto.String("." + fqn)

		fqn, dsc = l.resolve(fd, mtd.GetOutputType(), isMessage, scopes)
		if dsc == nil {
			return ErrorWithSourcePos{Pos: node.getOutputType().start(), Underlying: fmt.Errorf("%s: unknown response type %s", scope, mtd.GetOutputType())}
		}
		if _, ok := dsc.(*dpb.DescriptorProto); !ok {
			otherType := descriptorType(dsc)
			return ErrorWithSourcePos{Pos: node.getOutputType().start(), Underlying: fmt.Errorf("%s: invalid response type: %s is a %s, not a message", scope, fqn, otherType)}
		}
		mtd.OutputType = proto.String("." + fqn)
	}
	return nil
}

func (l *linker) resolveOptions(r *parseResult, fd *dpb.FileDescriptorProto, elemType, elemName, optType string, opts []*dpb.UninterpretedOption, scopes []scope) error {
	var scope string
	if elemType != "file" {
		scope = fmt.Sprintf("%s %s: ", elemType, elemName)
	}
	for _, opt := range opts {
		for _, nm := range opt.Name {
			if nm.GetIsExtension() {
				node := r.getOptionNamePartNode(nm)
				fqn, dsc := l.resolve(fd, nm.GetNamePart(), isField, scopes)
				if dsc == nil {
					return ErrorWithSourcePos{Pos: node.start(), Underlying: fmt.Errorf("%sunknown extension %s", scope, nm.GetNamePart())}
				}
				if ext, ok := dsc.(*dpb.FieldDescriptorProto); !ok {
					otherType := descriptorType(dsc)
					return ErrorWithSourcePos{Pos: node.start(), Underlying: fmt.Errorf("%sinvalid extension: %s is a %s, not an extension", scope, nm.GetNamePart(), otherType)}
				} else if ext.GetExtendee() == "" {
					return ErrorWithSourcePos{Pos: node.start(), Underlying: fmt.Errorf("%sinvalid extension: %s is a field but not an extension", scope, nm.GetNamePart())}
				}
				nm.NamePart = proto.String("." + fqn)
			}
		}
	}
	return nil
}

func (l *linker) resolve(fd *dpb.FileDescriptorProto, name string, allowed func(proto.Message) bool, scopes []scope) (string, proto.Message) {
	if strings.HasPrefix(name, ".") {
		// already fully-qualified
		d := l.findSymbol(fd, name[1:], false, map[*dpb.FileDescriptorProto]struct{}{})
		if d != nil {
			return name[1:], d
		}
	} else {
		// unqualified, so we look in the enclosing (last) scope first and move
		// towards outermost (first) scope, trying to resolve the symbol
		var bestGuess proto.Message
		var bestGuessFqn string
		for i := len(scopes) - 1; i >= 0; i-- {
			fqn, d := scopes[i](name)
			if d != nil {
				if allowed(d) {
					return fqn, d
				} else if bestGuess == nil {
					bestGuess = d
					bestGuessFqn = fqn
				}
			}
		}
		// we return best guess, even though it was not an allowed kind of
		// descriptor, so caller can print a better error message (e.g.
		// indicating that the name was found but that it's the wrong type)
		return bestGuessFqn, bestGuess
	}
	return "", nil
}

func isField(m proto.Message) bool {
	_, ok := m.(*dpb.FieldDescriptorProto)
	return ok
}

func isMessage(m proto.Message) bool {
	_, ok := m.(*dpb.DescriptorProto)
	return ok
}

func isType(m proto.Message) bool {
	switch m.(type) {
	case *dpb.DescriptorProto, *dpb.EnumDescriptorProto:
		return true
	}
	return false
}

// scope represents a lexical scope in a proto file in which messages and enums
// can be declared.
type scope func(string) (string, proto.Message)

func fileScope(fd *dpb.FileDescriptorProto, l *linker) scope {
	// we search symbols in this file, but also symbols in other files
	// that have the same package as this file
	pkg := fd.GetPackage()
	return func(name string) (string, proto.Message) {
		var n string
		if pkg == "" {
			n = name
		} else {
			n = pkg + "." + name
		}
		d := l.findSymbol(fd, n, false, map[*dpb.FileDescriptorProto]struct{}{})
		if d != nil {
			return n, d
		}
		// maybe name is already fully-qualified, just without a leading dot
		d = l.findSymbol(fd, name, false, map[*dpb.FileDescriptorProto]struct{}{})
		if d != nil {
			return name, d
		}
		return "", nil
	}
}

func messageScope(messageName string, filePool map[string]proto.Message) scope {
	return func(name string) (string, proto.Message) {
		n := messageName + "." + name
		if d, ok := filePool[n]; ok {
			return n, d
		}
		return "", nil
	}
}

func (l *linker) findSymbol(fd *dpb.FileDescriptorProto, name string, public bool, checked map[*dpb.FileDescriptorProto]struct{}) proto.Message {
	if _, ok := checked[fd]; ok {
		// already checked this one
		return nil
	}
	checked[fd] = struct{}{}
	d := l.descriptorPool[fd][name]
	if d != nil {
		return d
	}

	// When public = false, we are searching only directly imported symbols. But we
	// also need to search transitive public imports due to semantics of public imports.
	if public {
		for _, depIndex := range fd.PublicDependency {
			dep := fd.Dependency[depIndex]
			depres := l.files[dep]
			if depres == nil {
				// we'll catch this error later
				continue
			}
			if d = l.findSymbol(depres.fd, name, true, checked); d != nil {
				return d
			}
		}
	} else {
		for _, dep := range fd.Dependency {
			depres := l.files[dep]
			if depres == nil {
				// we'll catch this error later
				continue
			}
			if d = l.findSymbol(depres.fd, name, true, checked); d != nil {
				return d
			}
		}
	}

	return nil
}

func (l *linker) createdLinkedDescriptors() (map[string]*desc.FileDescriptor, error) {
	names := make([]string, 0, len(l.files))
	for name := range l.files {
		names = append(names, name)
	}
	sort.Strings(names)
	linked := map[string]*desc.FileDescriptor{}
	for _, name := range names {
		if _, err := l.linkFile(name, nil, linked); err != nil {
			return nil, err
		}
	}
	return linked, nil
}

func (l *linker) linkFile(name string, seen []string, linked map[string]*desc.FileDescriptor) (*desc.FileDescriptor, error) {
	// check for import cycle
	for _, s := range seen {
		if name == s {
			var msg bytes.Buffer
			first := true
			for _, s := range seen {
				if first {
					first = false
				} else {
					msg.WriteString(" -> ")
				}
				fmt.Fprintf(&msg, "%q", s)
			}
			fmt.Fprintf(&msg, " -> %q", name)
			return nil, fmt.Errorf("cycle found in imports: %s", msg.String())
		}
	}
	seen = append(seen, name)

	if lfd, ok := linked[name]; ok {
		// already linked
		return lfd, nil
	}
	r := l.files[name]
	if r == nil {
		importer := seen[len(seen)-2] // len-1 is *this* file, before that is the one that imported it
		return nil, fmt.Errorf("no descriptor found for %q, imported by %q", name, importer)
	}
	var deps []*desc.FileDescriptor
	for _, dep := range r.fd.Dependency {
		ldep, err := l.linkFile(dep, seen, linked)
		if err != nil {
			return nil, err
		}
		deps = append(deps, ldep)
	}
	lfd, err := desc.CreateFileDescriptor(r.fd, deps...)
	if err != nil {
		return nil, fmt.Errorf("error linking %q: %s", name, err)
	}
	linked[name] = lfd
	return lfd, nil
}

func (l *linker) interpretFileOptions(r *parseResult, fd *desc.FileDescriptor) error {
	opts := fd.GetFileOptions()
	if opts != nil {
		if len(opts.UninterpretedOption) > 0 {
			if err := l.interpretOptions(r, fd, opts, opts.UninterpretedOption); err != nil {
				return err
			}
		}
		opts.UninterpretedOption = nil
	}
	for _, md := range fd.GetMessageTypes() {
		if err := l.interpretMessageOptions(r, md); err != nil {
			return err
		}
	}
	for _, fld := range fd.GetExtensions() {
		if err := l.interpretFieldOptions(r, fld); err != nil {
			return err
		}
	}
	for _, ed := range fd.GetEnumTypes() {
		if err := l.interpretEnumOptions(r, ed); err != nil {
			return err
		}
	}
	for _, sd := range fd.GetServices() {
		opts := sd.GetServiceOptions()
		if opts != nil {
			if len(opts.UninterpretedOption) > 0 {
				if err := l.interpretOptions(r, sd, opts, opts.UninterpretedOption); err != nil {
					return err
				}
			}
			opts.UninterpretedOption = nil
		}
		for _, mtd := range sd.GetMethods() {
			opts := mtd.GetMethodOptions()
			if opts != nil {
				if len(opts.UninterpretedOption) > 0 {
					if err := l.interpretOptions(r, mtd, opts, opts.UninterpretedOption); err != nil {
						return err
					}
				}
				opts.UninterpretedOption = nil
			}
		}
	}
	return nil
}

func (l *linker) interpretMessageOptions(r *parseResult, md *desc.MessageDescriptor) error {
	opts := md.GetMessageOptions()
	if opts != nil {
		if len(opts.UninterpretedOption) > 0 {
			if err := l.interpretOptions(r, md, opts, opts.UninterpretedOption); err != nil {
				return err
			}
		}
		opts.UninterpretedOption = nil
	}
	for _, fld := range md.GetFields() {
		if err := l.interpretFieldOptions(r, fld); err != nil {
			return err
		}
	}
	for _, fld := range md.GetNestedExtensions() {
		if err := l.interpretFieldOptions(r, fld); err != nil {
			return err
		}
	}
	for _, er := range md.AsDescriptorProto().GetExtensionRange() {
		opts := er.Options
		if opts != nil {
			if len(opts.UninterpretedOption) > 0 {
				d := extRangeDescriptorish{md: md, er: er}
				if err := l.interpretOptions(r, d, opts, opts.UninterpretedOption); err != nil {
					return err
				}
			}
			opts.UninterpretedOption = nil
		}
	}
	for _, nmd := range md.GetNestedMessageTypes() {
		if err := l.interpretMessageOptions(r, nmd); err != nil {
			return err
		}
	}
	for _, ed := range md.GetNestedEnumTypes() {
		if err := l.interpretEnumOptions(r, ed); err != nil {
			return err
		}
	}
	return nil
}

type extRangeDescriptorish struct {
	md *desc.MessageDescriptor
	er *dpb.DescriptorProto_ExtensionRange
}

func (er extRangeDescriptorish) GetFile() *desc.FileDescriptor {
	return er.md.GetFile()
}

func (er extRangeDescriptorish) GetName() string {
	return fmt.Sprintf("%s:%d-%d", er.md.GetName(), er.er.GetStart(), er.er.GetEnd()-1)
}

func (er extRangeDescriptorish) AsProto() proto.Message {
	return er.er
}

func (l *linker) interpretFieldOptions(r *parseResult, fld *desc.FieldDescriptor) error {
	opts := fld.GetFieldOptions()
	if opts != nil {
		if len(opts.UninterpretedOption) > 0 {
			uo := opts.UninterpretedOption
			scope := fmt.Sprintf("field %s", fld.GetFullyQualifiedName())

			// process json_name pseudo-option
			if index, err := findOption(r, scope, uo, "json_name"); err != nil {
				return err
			} else if index >= 0 {
				opt := uo[index]
				optNode := r.getOptionNode(opt)

				// attribute source code info
				if on, ok := optNode.(*optionNode); ok {
					r.interpretedOptions[on] = []int32{-1, internal.Field_jsonNameTag}
				}
				uo = removeOption(uo, index)
				if opt.StringValue == nil {
					return ErrorWithSourcePos{Pos: optNode.getValue().start(), Underlying: fmt.Errorf("%s: expecting string value for json_name option", scope)}
				}
				fld.AsFieldDescriptorProto().JsonName = proto.String(string(opt.StringValue))
			}

			// and process default pseudo-option
			if i, err := processDefaultOption(r, scope, fld, uo); err != nil {
				return err
			} else if i >= 0 {
				// attribute source code info
				optNode := r.getOptionNode(uo[i])
				if on, ok := optNode.(*optionNode); ok {
					r.interpretedOptions[on] = []int32{-1, internal.Field_defaultTag}
				}
				uo = removeOption(uo, i)
			}

			if len(uo) == 0 {
				// no real options, only pseudo-options above? clear out options
				fld.AsFieldDescriptorProto().Options = nil
			} else if err := l.interpretOptions(r, fld, opts, uo); err != nil {
				return err
			}
		}
		opts.UninterpretedOption = nil
	}
	return nil
}

func processDefaultOption(res *parseResult, scope string, fld *desc.FieldDescriptor, uos []*dpb.UninterpretedOption) (defaultIndex int, err error) {
	found, err := findOption(res, scope, uos, "default")
	if err != nil {
		return -1, err
	} else if found == -1 {
		return -1, nil
	}
	opt := uos[found]
	optNode := res.getOptionNode(opt)
	if fld.IsRepeated() {
		return -1, ErrorWithSourcePos{Pos: optNode.getName().start(), Underlying: fmt.Errorf("%s: default value cannot be set because field is repeated", scope)}
	}
	if fld.GetType() == dpb.FieldDescriptorProto_TYPE_GROUP || fld.GetType() == dpb.FieldDescriptorProto_TYPE_MESSAGE {
		return -1, ErrorWithSourcePos{Pos: optNode.getName().start(), Underlying: fmt.Errorf("%s: default value cannot be set because field is a message", scope)}
	}
	val := optNode.getValue()
	if _, ok := val.(*aggregateLiteralNode); ok {
		return -1, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%s: default value cannot be an aggregate", scope)}
	}
	mc := &messageContext{
		res:         res,
		file:        fld.GetFile(),
		elementName: fld.GetFullyQualifiedName(),
		elementType: descriptorType(fld.AsProto()),
		option:      opt,
	}
	v, err := fieldValue(res, mc, fld, val, true)
	if err != nil {
		return -1, err
	}
	if str, ok := v.(string); ok {
		fld.AsFieldDescriptorProto().DefaultValue = proto.String(str)
	} else if b, ok := v.([]byte); ok {
		fld.AsFieldDescriptorProto().DefaultValue = proto.String(encodeDefaultBytes(b))
	} else {
		var flt float64
		var ok bool
		if flt, ok = v.(float64); !ok {
			var flt32 float32
			if flt32, ok = v.(float32); ok {
				flt = float64(flt32)
			}
		}
		if ok {
			if math.IsInf(flt, 1) {
				fld.AsFieldDescriptorProto().DefaultValue = proto.String("inf")
			} else if ok && math.IsInf(flt, -1) {
				fld.AsFieldDescriptorProto().DefaultValue = proto.String("-inf")
			} else if ok && math.IsNaN(flt) {
				fld.AsFieldDescriptorProto().DefaultValue = proto.String("nan")
			} else {
				fld.AsFieldDescriptorProto().DefaultValue = proto.String(fmt.Sprintf("%v", v))
			}
		} else {
			fld.AsFieldDescriptorProto().DefaultValue = proto.String(fmt.Sprintf("%v", v))
		}
	}
	return found, nil
}

func encodeDefaultBytes(b []byte) string {
	var buf bytes.Buffer
	writeEscapedBytes(&buf, b)
	return buf.String()
}

func (l *linker) interpretEnumOptions(r *parseResult, ed *desc.EnumDescriptor) error {
	opts := ed.GetEnumOptions()
	if opts != nil {
		if len(opts.UninterpretedOption) > 0 {
			if err := l.interpretOptions(r, ed, opts, opts.UninterpretedOption); err != nil {
				return err
			}
		}
		opts.UninterpretedOption = nil
	}
	for _, evd := range ed.GetValues() {
		opts := evd.GetEnumValueOptions()
		if opts != nil {
			if len(opts.UninterpretedOption) > 0 {
				if err := l.interpretOptions(r, evd, opts, opts.UninterpretedOption); err != nil {
					return err
				}
			}
			opts.UninterpretedOption = nil
		}
	}
	return nil
}

type descriptorish interface {
	GetFile() *desc.FileDescriptor
	GetName() string
	AsProto() proto.Message
}

func (l *linker) interpretOptions(res *parseResult, element descriptorish, opts proto.Message, uninterpreted []*dpb.UninterpretedOption) error {
	optsd, err := desc.LoadMessageDescriptorForMessage(opts)
	if err != nil {
		return err
	}
	dm := dynamic.NewMessage(optsd)
	err = dm.ConvertFrom(opts)
	if err != nil {
		node := res.nodes[element.AsProto()]
		return ErrorWithSourcePos{Pos: node.start(), Underlying: err}
	}

	mc := &messageContext{res: res, file: element.GetFile(), elementName: element.GetName(), elementType: descriptorType(element.AsProto())}
	for _, uo := range uninterpreted {
		node := res.getOptionNode(uo)
		if !uo.Name[0].GetIsExtension() && uo.Name[0].GetNamePart() == "uninterpreted_option" {
			// uninterpreted_option might be found reflectively, but is not actually valid for use
			return ErrorWithSourcePos{Pos: node.getName().start(), Underlying: fmt.Errorf("%vinvalid option 'uninterpreted_option'", mc)}
		}
		mc.option = uo
		path, err := l.interpretField(res, mc, element, dm, uo, 0, nil)
		if err != nil {
			return err
		}
		if optn, ok := node.(*optionNode); ok {
			res.interpretedOptions[optn] = path
		}
	}

	if err := dm.ValidateRecursive(); err != nil {
		node := res.nodes[element.AsProto()]
		return ErrorWithSourcePos{Pos: node.start(), Underlying: fmt.Errorf("error in %s options: %v", descriptorType(element.AsProto()), err)}
	}

	if err := dm.ConvertTo(opts); err != nil {
		node := res.nodes[element.AsProto()]
		return ErrorWithSourcePos{Pos: node.start(), Underlying: err}
	}
	return nil
}

func (l *linker) interpretField(res *parseResult, mc *messageContext, element descriptorish, dm *dynamic.Message, opt *dpb.UninterpretedOption, nameIndex int, pathPrefix []int32) (path []int32, err error) {
	var fld *desc.FieldDescriptor
	nm := opt.GetName()[nameIndex]
	node := res.getOptionNamePartNode(nm)
	if nm.GetIsExtension() {
		extName := nm.GetNamePart()[1:] /* skip leading dot */
		fld = findExtension(element.GetFile(), extName, false, map[*desc.FileDescriptor]struct{}{})
		if fld == nil {
			return nil, ErrorWithSourcePos{
				Pos: node.start(),
				Underlying: fmt.Errorf("%vunrecognized extension %s of %s",
					mc, extName, dm.GetMessageDescriptor().GetFullyQualifiedName()),
			}
		}
		if fld.GetOwner().GetFullyQualifiedName() != dm.GetMessageDescriptor().GetFullyQualifiedName() {
			return nil, ErrorWithSourcePos{
				Pos: node.start(),
				Underlying: fmt.Errorf("%vextension %s should extend %s but instead extends %s",
					mc, extName, dm.GetMessageDescriptor().GetFullyQualifiedName(), fld.GetOwner().GetFullyQualifiedName()),
			}
		}
	} else {
		fld = dm.GetMessageDescriptor().FindFieldByName(nm.GetNamePart())
		if fld == nil {
			return nil, ErrorWithSourcePos{
				Pos: node.start(),
				Underlying: fmt.Errorf("%vfield %s of %s does not exist",
					mc, nm.GetNamePart(), dm.GetMessageDescriptor().GetFullyQualifiedName()),
			}
		}
	}

	path = append(pathPrefix, fld.GetNumber())

	if len(opt.GetName()) > nameIndex+1 {
		nextnm := opt.GetName()[nameIndex+1]
		nextnode := res.getOptionNamePartNode(nextnm)
		if fld.GetType() != dpb.FieldDescriptorProto_TYPE_MESSAGE {
			return nil, ErrorWithSourcePos{
				Pos: nextnode.start(),
				Underlying: fmt.Errorf("%vcannot set field %s because %s is not a message",
					mc, nextnm.GetNamePart(), nm.GetNamePart()),
			}
		}
		if fld.IsRepeated() {
			return nil, ErrorWithSourcePos{
				Pos: nextnode.start(),
				Underlying: fmt.Errorf("%vcannot set field %s because %s is repeated (must use an aggregate)",
					mc, nextnm.GetNamePart(), nm.GetNamePart()),
			}
		}
		var fdm *dynamic.Message
		var err error
		if dm.HasField(fld) {
			var v interface{}
			v, err = dm.TryGetField(fld)
			fdm, _ = v.(*dynamic.Message)
		} else {
			fdm = dynamic.NewMessage(fld.GetMessageType())
			err = dm.TrySetField(fld, fdm)
		}
		if err != nil {
			return nil, ErrorWithSourcePos{Pos: node.start(), Underlying: err}
		}
		// recurse to set next part of name
		return l.interpretField(res, mc, element, fdm, opt, nameIndex+1, path)
	}

	optNode := res.getOptionNode(opt)
	if err := setOptionField(res, mc, dm, fld, node, optNode.getValue()); err != nil {
		return nil, err
	}
	if fld.IsRepeated() {
		path = append(path, int32(dm.FieldLength(fld))-1)
	}
	return path, nil
}

func findExtension(fd *desc.FileDescriptor, name string, public bool, checked map[*desc.FileDescriptor]struct{}) *desc.FieldDescriptor {
	if _, ok := checked[fd]; ok {
		return nil
	}
	checked[fd] = struct{}{}
	d := fd.FindSymbol(name)
	if d != nil {
		if fld, ok := d.(*desc.FieldDescriptor); ok {
			return fld
		}
		return nil
	}

	// When public = false, we are searching only directly imported symbols. But we
	// also need to search transitive public imports due to semantics of public imports.
	if public {
		for _, dep := range fd.GetPublicDependencies() {
			d := findExtension(dep, name, true, checked)
			if d != nil {
				return d
			}
		}
	} else {
		for _, dep := range fd.GetDependencies() {
			d := findExtension(dep, name, true, checked)
			if d != nil {
				return d
			}
		}
	}
	return nil
}

func setOptionField(res *parseResult, mc *messageContext, dm *dynamic.Message, fld *desc.FieldDescriptor, name node, val valueNode) error {
	v := val.value()
	if sl, ok := v.([]valueNode); ok {
		// handle slices a little differently than the others
		if !fld.IsRepeated() {
			return ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue is an array but field is not repeated", mc)}
		}
		origPath := mc.optAggPath
		defer func() {
			mc.optAggPath = origPath
		}()
		for index, item := range sl {
			mc.optAggPath = fmt.Sprintf("%s[%d]", origPath, index)
			if v, err := fieldValue(res, mc, fld, item, false); err != nil {
				return err
			} else if err = dm.TryAddRepeatedField(fld, v); err != nil {
				return ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%verror setting value: %s", mc, err)}
			}
		}
		return nil
	}

	v, err := fieldValue(res, mc, fld, val, false)
	if err != nil {
		return err
	}
	if fld.IsRepeated() {
		err = dm.TryAddRepeatedField(fld, v)
	} else {
		if dm.HasField(fld) {
			return ErrorWithSourcePos{Pos: name.start(), Underlying: fmt.Errorf("%vnon-repeated option field %s already set", mc, fieldName(fld))}
		}
		err = dm.TrySetField(fld, v)
	}
	if err != nil {
		return ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%verror setting value: %s", mc, err)}
	}

	return nil
}

type messageContext struct {
	res         *parseResult
	file        *desc.FileDescriptor
	elementType string
	elementName string
	option      *dpb.UninterpretedOption
	optAggPath  string
}

func (c *messageContext) String() string {
	var ctx bytes.Buffer
	if c.elementType != "file" {
		fmt.Fprintf(&ctx, "%s %s: ", c.elementType, c.elementName)
	}
	if c.option != nil && c.option.Name != nil {
		ctx.WriteString("option ")
		writeOptionName(&ctx, c.option.Name)
		if c.res.nodes == nil {
			// if we have no source position info, try to provide as much context
			// as possible (if nodes != nil, we don't need this because any errors
			// will actually have file and line numbers)
			if c.optAggPath != "" {
				fmt.Fprintf(&ctx, " at %s", c.optAggPath)
			}
		}
		ctx.WriteString(": ")
	}
	return ctx.String()
}

func writeOptionName(buf *bytes.Buffer, parts []*dpb.UninterpretedOption_NamePart) {
	first := true
	for _, p := range parts {
		if first {
			first = false
		} else {
			buf.WriteByte('.')
		}
		nm := p.GetNamePart()
		if nm[0] == '.' {
			// skip leading dot
			nm = nm[1:]
		}
		if p.GetIsExtension() {
			buf.WriteByte('(')
			buf.WriteString(nm)
			buf.WriteByte(')')
		} else {
			buf.WriteString(nm)
		}
	}
}

func fieldName(fld *desc.FieldDescriptor) string {
	if fld.IsExtension() {
		return fld.GetFullyQualifiedName()
	} else {
		return fld.GetName()
	}
}

func valueKind(val interface{}) string {
	switch val := val.(type) {
	case identifier:
		return "identifier"
	case bool:
		return "bool"
	case int64:
		if val < 0 {
			return "negative integer"
		}
		return "integer"
	case uint64:
		return "integer"
	case float64:
		return "double"
	case string, []byte:
		return "string"
	case []*aggregateEntryNode:
		return "message"
	default:
		return fmt.Sprintf("%T", val)
	}
}

func fieldValue(res *parseResult, mc *messageContext, fld *desc.FieldDescriptor, val valueNode, enumAsString bool) (interface{}, error) {
	v := val.value()
	switch fld.GetType() {
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		if id, ok := v.(identifier); ok {
			ev := fld.GetEnumType().FindValueByName(string(id))
			if ev == nil {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%venum %s has no value named %s", mc, fld.GetEnumType().GetFullyQualifiedName(), id)}
			}
			if enumAsString {
				return ev.GetName(), nil
			} else {
				return ev.GetNumber(), nil
			}
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting enum, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_MESSAGE, dpb.FieldDescriptorProto_TYPE_GROUP:
		if aggs, ok := v.([]*aggregateEntryNode); ok {
			fmd := fld.GetMessageType()
			fdm := dynamic.NewMessage(fmd)
			origPath := mc.optAggPath
			defer func() {
				mc.optAggPath = origPath
			}()
			for _, a := range aggs {
				if origPath == "" {
					mc.optAggPath = a.name.value()
				} else {
					mc.optAggPath = origPath + "." + a.name.value()
				}
				var ffld *desc.FieldDescriptor
				if a.name.isExtension {
					n := a.name.name.val
					ffld = findExtension(mc.file, n, false, map[*desc.FileDescriptor]struct{}{})
					if ffld == nil {
						// may need to qualify with package name
						pkg := mc.file.GetPackage()
						if pkg != "" {
							ffld = findExtension(mc.file, pkg+"."+n, false, map[*desc.FileDescriptor]struct{}{})
						}
					}
				} else {
					ffld = fmd.FindFieldByName(a.name.value())
				}
				if ffld == nil {
					return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vfield %s not found", mc, a.name.name.val)}
				}
				if err := setOptionField(res, mc, fdm, ffld, a.name, a.val); err != nil {
					return nil, err
				}
			}
			return fdm, nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting message, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_BOOL:
		if b, ok := v.(bool); ok {
			return b, nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting bool, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		if str, ok := v.(string); ok {
			return []byte(str), nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting bytes, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_STRING:
		if str, ok := v.(string); ok {
			return str, nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting string, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_INT32, dpb.FieldDescriptorProto_TYPE_SINT32, dpb.FieldDescriptorProto_TYPE_SFIXED32:
		if i, ok := v.(int64); ok {
			if i > math.MaxInt32 || i < math.MinInt32 {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %d is out of range for int32", mc, i)}
			}
			return int32(i), nil
		}
		if ui, ok := v.(uint64); ok {
			if ui > math.MaxInt32 {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %d is out of range for int32", mc, ui)}
			}
			return int32(ui), nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting int32, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_UINT32, dpb.FieldDescriptorProto_TYPE_FIXED32:
		if i, ok := v.(int64); ok {
			if i > math.MaxUint32 || i < 0 {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %d is out of range for uint32", mc, i)}
			}
			return uint32(i), nil
		}
		if ui, ok := v.(uint64); ok {
			if ui > math.MaxUint32 {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %d is out of range for uint32", mc, ui)}
			}
			return uint32(ui), nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting uint32, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_INT64, dpb.FieldDescriptorProto_TYPE_SINT64, dpb.FieldDescriptorProto_TYPE_SFIXED64:
		if i, ok := v.(int64); ok {
			return i, nil
		}
		if ui, ok := v.(uint64); ok {
			if ui > math.MaxInt64 {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %d is out of range for int64", mc, ui)}
			}
			return int64(ui), nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting int64, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_UINT64, dpb.FieldDescriptorProto_TYPE_FIXED64:
		if i, ok := v.(int64); ok {
			if i < 0 {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %d is out of range for uint64", mc, i)}
			}
			return uint64(i), nil
		}
		if ui, ok := v.(uint64); ok {
			return ui, nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting uint64, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_DOUBLE:
		if d, ok := v.(float64); ok {
			return d, nil
		}
		if i, ok := v.(int64); ok {
			return float64(i), nil
		}
		if u, ok := v.(uint64); ok {
			return float64(u), nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting double, got %s", mc, valueKind(v))}
	case dpb.FieldDescriptorProto_TYPE_FLOAT:
		if d, ok := v.(float64); ok {
			if (d > math.MaxFloat32 || d < -math.MaxFloat32) && !math.IsInf(d, 1) && !math.IsInf(d, -1) && !math.IsNaN(d) {
				return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vvalue %f is out of range for float", mc, d)}
			}
			return float32(d), nil
		}
		if i, ok := v.(int64); ok {
			return float32(i), nil
		}
		if u, ok := v.(uint64); ok {
			return float32(u), nil
		}
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vexpecting float, got %s", mc, valueKind(v))}
	default:
		return nil, ErrorWithSourcePos{Pos: val.start(), Underlying: fmt.Errorf("%vunrecognized field type: %s", mc, fld.GetType())}
	}
}
