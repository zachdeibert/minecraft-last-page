package gpu

//	#cgo pkg-config: OpenCL
//	#cgo CFLAGS: -Wall
//	#include <stdlib.h>
//	#include <CL/cl.h>
//
//	static void cleanup_queue_early(void *_queue) {
//		cl_command_queue *queue = (cl_command_queue *) _queue;
//		clFlush(*queue);
//		clFinish(*queue);
//	}
import "C"

func init() {
	cleanupFuncs = append(cleanupFuncs, func() {
		C.cleanup_queue_early(queue)
	})
}
