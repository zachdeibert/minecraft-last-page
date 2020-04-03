package com.github.zachdeibert.minecraftlastpage.sanitycheck

import com.google.common.base.Stopwatch
import java.math.BigInteger

private const val MULTITASKING_THRESHOLD: Int = 3
private const val NUM_WORKERS: Int = 16
private const val WORKER_TIME: Int = 10
private const val MAX_STRING_LENGTH: Int = 6

private fun createSieves(sieves: MutableList<Sieve>, prefix: String) {
	for (c in Sieve.ALPHABET) {
		val s = prefix + c
		if (s.length == MULTITASKING_THRESHOLD) {
			sieves.add(Sieve(s, MAX_STRING_LENGTH))
		} else {
			createSieves(sieves, s)
		}
	}
}

fun main(args: Array<String>) {
	val sieves = mutableListOf<Sieve>(
			Sieve("", MULTITASKING_THRESHOLD)
	)
	createSieves(sieves, "")
	val pool = Pool(sieves)
	val watch = Stopwatch.createStarted()
	val workers = Array<Worker>(NUM_WORKERS) {
		Worker(pool, WORKER_TIME).apply {
			start()
		}
	}
	System.out.println("Worker threads started.")
	for (worker in workers) {
		worker.join()
	}
	watch.stop()
	var num = BigInteger.ZERO
	for (sieve in sieves) {
		num += BigInteger.valueOf(sieve.count)
	}
	System.out.printf("Done with %s hashes in %s.\n", num, watch.elapsed())
}
