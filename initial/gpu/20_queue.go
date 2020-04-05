package gpu

//	#cgo pkg-config: OpenCL
//	#cgo CFLAGS: -Wall
//	#include <errno.h>
//	#include <stdlib.h>
//	#include <CL/cl.h>
//
//	static void *create_queue(void *_ctx, void *_device_id) {
//		cl_context *ctx = (cl_context *) _ctx;
//		cl_device_id *device_id = (cl_device_id *) _device_id;
//		cl_command_queue *queue = (cl_command_queue *) malloc(sizeof(cl_command_queue));
//		if (queue == NULL) {
//			errno = ENOBUFS;
//			return NULL;
//		}
//		cl_int ret;
//		*queue = clCreateCommandQueue(*ctx, *device_id, 0, &ret);
//		if (ret != CL_SUCCESS) {
//			free(queue);
//			errno = ret;
//			return NULL;
//		}
//		errno = 0;
//		return queue;
//	}
//
//	static void write_buffer(void *_queue, void *_buffer, int offset, int length, void *data) {
//		cl_command_queue *queue = (cl_command_queue *) _queue;
//		cl_mem *buffer = (cl_mem *) _buffer;
//		cl_int ret;
//		if ((ret = clEnqueueWriteBuffer(*queue, *buffer, CL_TRUE, offset, length, data, 0, NULL, NULL)) != CL_SUCCESS) {
//			errno = ret;
//			return;
//		}
//		errno = 0;
//	}
//
//	static void read_buffer(void *_queue, void *_buffer, int offset, int length, void *data) {
//		cl_command_queue *queue = (cl_command_queue *) _queue;
//		cl_mem *buffer = (cl_mem *) _buffer;
//		cl_int ret;
//		if ((ret = clEnqueueReadBuffer(*queue, *buffer, CL_TRUE, offset, length, data, 0, NULL, NULL)) != CL_SUCCESS) {
//			errno = ret;
//			return;
//		}
//		errno = 0;
//	}
//
//	static void run_kernel(void *_queue, void *_kernel, int work_dim, int gbl1, int gbl2, int gbl3) {
//		cl_command_queue *queue = (cl_command_queue *) _queue;
//		cl_kernel *kernel = (cl_kernel *) _kernel;
//		size_t gbl[3] = { gbl1, gbl2, gbl3 };
//		cl_int ret;
//		if ((ret = clEnqueueNDRangeKernel(*queue, *kernel, work_dim, NULL, gbl, NULL, 0, NULL, NULL)) != CL_SUCCESS) {
//			errno = ret;
//			return;
//		}
//		if ((ret = clFlush(*queue)) != CL_SUCCESS) {
//			errno = ret;
//			return;
//		}
//		if ((ret = clFinish(*queue)) != CL_SUCCESS) {
//			errno = ret;
//			return;
//		}
//		errno = 0;
//	}
//
//	static void cleanup_queue(void *_queue) {
//		cl_command_queue *queue = (cl_command_queue *) _queue;
//		clReleaseCommandQueue(*queue);
//		free(queue);
//	}
import "C"
import (
	"errors"
	"unsafe"
)

var (
	queue unsafe.Pointer
)

func writeBuffer(buf buffer, offset int, length int, data []byte) error {
	_, e := C.write_buffer(queue, buf.ptr, C.int(offset), C.int(length), unsafe.Pointer(&data[0]))
	return convertError(e)
}

func readBuffer(buf buffer, offset int, length int, data []byte) error {
	_, e := C.read_buffer(queue, buf.ptr, C.int(offset), C.int(length), unsafe.Pointer(&data[0]))
	return convertError(e)
}

func runKernel(kernel kernel, global []int) error {
	if len(global) == 0 || len(global) > 3 {
		return errors.New("Invalid work group size")
	}
	l := len(global)
	global = append(global, 0, 0)
	_, e := C.run_kernel(queue, kernel.ptr, C.int(l), C.int(global[0]), C.int(global[1]), C.int(global[2]))
	return convertError(e)
}

func init() {
	var e error
	queue, e = C.create_queue(context, deviceID)
	if e != nil {
		panic(convertError(e))
	}
	cleanupFuncs = append(cleanupFuncs, func() {
		C.cleanup_queue(queue)
	})
}
