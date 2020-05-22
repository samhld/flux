package libflux

// #include "influxdata/flux.h"
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/influxdata/flux/codes"
	"github.com/influxdata/flux/internal/errors"
)

// SemanticPkg is a Rust pointer to a semantic package.
type SemanticPkg struct {
	ptr  *C.struct_flux_semantic_pkg_t
	free func(ptr *C.struct_flux_semantic_pkg_t)
}

// MarshalFB serializes the given semantic package into a flatbuffer.
func (p *SemanticPkg) MarshalFB() ([]byte, error) {
	var buf C.struct_flux_buffer_t
	if err := C.flux_semantic_marshal_fb(p.ptr, &buf); err != nil {
		defer C.flux_free_error(err)
		cstr := C.flux_error_str(err)
		defer C.flux_free_bytes(cstr)

		str := C.GoString(cstr)
		return nil, errors.Newf(codes.Internal, "could not marshal semantic graph to FlatBuffer: %v", str)
	}
	// See MarshalFB in ASTPkg for why this is needed.
	runtime.KeepAlive(p)
	defer C.flux_free_bytes(buf.data)

	data := C.GoBytes(unsafe.Pointer(buf.data), C.int(buf.len))
	return data, nil
}

// Free frees the memory allocated by Rust for the semantic graph.
func (p *SemanticPkg) Free() {
	if p.ptr != nil {
		if p.free != nil {
			p.free(p.ptr)
		} else {
			C.flux_free_semantic_pkg(p.ptr)
		}
	}
	p.ptr = nil

	// See the equivalent method in ASTPkg for why
	// this is needed.
	runtime.KeepAlive(p)
}

// Analyze parses the given Flux source, performs type inference
// (taking into account types from prelude and stdlib) and returns
// an a SemanticPkg containing an opaque pointer to the semantic graph.
// The graph can be deserialized by calling MarshalFB.
//
// Note that Analyze will consume the AST, so astPkg.ptr will be set to nil,
// even if there's an error in analysis.
func Analyze(astPkg *ASTPkg) (*SemanticPkg, error) {
	var semPkg *C.struct_flux_semantic_pkg_t
	defer func() {
		astPkg.ptr = nil
	}()
	if err := C.flux_analyze(astPkg.ptr, &semPkg); err != nil {
		defer C.flux_free_error(err)
		cstr := C.flux_error_str(err)
		defer C.flux_free_bytes(cstr)

		str := C.GoString(cstr)
		return nil, errors.New(codes.Invalid, str)
	}
	runtime.KeepAlive(astPkg)
	p := &SemanticPkg{ptr: semPkg}
	runtime.SetFinalizer(p, free)
	return p, nil
}

func AnalyzePackage(pkgpath string, astPkg *ASTPkg, imports *SemanticPkgSet) (*SemanticPkg, error) {
	defer func() { astPkg.ptr = nil }()

	cpkgpath := C.CString(pkgpath)
	defer C.free(unsafe.Pointer(cpkgpath))

	var importsPtr *C.struct_flux_semantic_pkgset_t
	if imports != nil {
		importsPtr = imports.ptr
	}
	semPkg := C.flux_analyze_package(cpkgpath, astPkg.ptr, importsPtr)
	if err := C.flux_semantic_pkg_get_error(semPkg); err != nil {
		C.flux_semantic_pkg_free(semPkg)
		defer C.flux_free_error(err)
		cstr := C.flux_error_str(err)
		defer C.flux_free_bytes(cstr)

		str := C.GoString(cstr)
		return nil, errors.New(codes.Invalid, str)
	}
	runtime.KeepAlive(astPkg)
	runtime.KeepAlive(imports)
	p := &SemanticPkg{
		ptr: semPkg,
		free: func(ptr *C.struct_flux_semantic_pkg_t) {
			C.flux_semantic_pkg_free(ptr)
		},
	}
	runtime.SetFinalizer(p, free)
	return p, nil
}

type SemanticPkgSet struct {
	ptr *C.struct_flux_semantic_pkgset_t
}

func NewSemanticPkgSet() *SemanticPkgSet {
	ptr := C.flux_semantic_pkgset_new()
	p := &SemanticPkgSet{ptr: ptr}
	runtime.SetFinalizer(p, free)
	return p
}

func (p *SemanticPkgSet) Add(pkg *SemanticPkg) error {
	defer func() { pkg.ptr = nil }()

	if err := C.flux_semantic_pkgset_add(p.ptr, pkg.ptr); err != nil {
		defer C.flux_free_error(err)
		cstr := C.flux_error_str(err)
		defer C.flux_free_bytes(cstr)

		str := C.GoString(cstr)
		return errors.New(codes.Invalid, str)
	}
	runtime.KeepAlive(pkg)
	runtime.KeepAlive(p)
	return nil
}

func (p *SemanticPkgSet) Free() {
	if p.ptr != nil {
		C.flux_semantic_pkgset_free(p.ptr)
	}
	p.ptr = nil

	// See the equivalent method in ASTPkg for why
	// this is needed.
	runtime.KeepAlive(p)
}

type Analyzer struct {
	ptr *C.struct_flux_semantic_analyzer_t
}

func NewAnalyzer(pkgpath string) *Analyzer {
	cstr := C.CString(pkgpath)
	defer C.free(unsafe.Pointer(cstr))

	ptr := C.flux_new_semantic_analyzer(cstr)
	p := &Analyzer{ptr: ptr}
	runtime.SetFinalizer(p, free)
	return p
}

func (p *Analyzer) Analyze(astPkg *ASTPkg) (*SemanticPkg, error) {
	var semPkg *C.struct_flux_semantic_pkg_t
	defer func() {
		astPkg.ptr = nil
	}()
	if err := C.flux_analyze_with(p.ptr, astPkg.ptr, &semPkg); err != nil {
		defer C.flux_free_error(err)
		cstr := C.flux_error_str(err)
		defer C.flux_free_bytes(cstr)

		str := C.GoString(cstr)
		return nil, errors.New(codes.Invalid, str)
	}
	runtime.KeepAlive(p)

	pkg := &SemanticPkg{ptr: semPkg}
	runtime.SetFinalizer(pkg, free)
	return pkg, nil
}

// Free frees the memory allocated by Rust for the semantic graph.
func (p *Analyzer) Free() {
	if p.ptr != nil {
		C.flux_free_semantic_analyzer(p.ptr)
	}
	p.ptr = nil

	// See the equivalent method in ASTPkg for why
	// this is needed.
	runtime.KeepAlive(p)
}

// EnvStdlib takes care of creating a flux_buffer_t, passes the buffer to
// the Flatbuffers TypeEnvironment and then takes care of freeing the data
func EnvStdlib() []byte {
	var buf C.struct_flux_buffer_t
	C.flux_get_env_stdlib(&buf)
	defer C.flux_free_bytes(buf.data)
	return C.GoBytes(unsafe.Pointer(buf.data), C.int(buf.len))
}
