package gpu

//	#cgo pkg-config: OpenCL
//	#cgo CFLAGS: -Wall
//	#include <errno.h>
//	#include <stdlib.h>
//	#include <CL/cl.h>
//
//	static void *create_buffer(void *_ctx, int mode, int size) {
//		cl_context *ctx = (cl_context *) _ctx;
//		cl_mem_flags flags;
//		switch (mode) {
//			case 0:
//				flags = CL_MEM_READ_WRITE;
//				break;
//			case 1:
//				flags = CL_MEM_WRITE_ONLY;
//				break;
//			case 2:
//				flags = CL_MEM_READ_ONLY;
//				break;
//			default:
//				return NULL;
//		}
//		cl_mem *mem = (cl_mem *) malloc(sizeof(cl_mem));
//		if (mem == NULL) {
//			errno = ENOBUFS;
//			return NULL;
//		}
//		cl_int ret;
//		*mem = clCreateBuffer(*ctx, flags, size, NULL, &ret);
//		if (ret != CL_SUCCESS) {
//			free(mem);
//			errno = ret;
//			return NULL;
//		}
//		errno = 0;
//		return mem;
//	}
//
//	static void cleanup_buffer(void *_mem) {
//		cl_mem *mem = (cl_mem *) _mem;
//		clReleaseMemObject(*mem);
//		free(mem);
//	}
import "C"
import (
	"unsafe"
)

type buffer struct {
	ptr unsafe.Pointer
}
type memMode int

const (
	memReadWrite memMode = 0
	memWriteOnly memMode = 1
	memReadOnly  memMode = 2
)

func createBuffer(mode memMode, size int) (buffer, error) {
	b, e := C.create_buffer(context, C.int(mode), C.int(size))
	return buffer{
		ptr: b,
	}, convertError(e)
}

func (b buffer) Close() {
	C.cleanup_buffer(b.ptr)
}
