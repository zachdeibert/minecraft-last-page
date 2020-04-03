package com.github.zachdeibert.minecraftlastpage.sanitycheck

import java.util.concurrent.locks.ReentrantLock

class Pool {
	private val sieves: Array<Sieve?>
	private var obtainIdx: Int
	private var releaseIdx: Int
	private val lock: ReentrantLock
	
	fun obtain(): Sieve? {
		try {
			lock.lock()
			if (obtainIdx == releaseIdx) {
				return null
			}
			val s = sieves[obtainIdx]
			obtainIdx = (obtainIdx + 1) % sieves.size
			return s
		} finally {
			lock.unlock()
		}
	}
	
	fun release(sieve: Sieve) {
		try {
			lock.lock()
			sieves[releaseIdx] = sieve
			releaseIdx = (releaseIdx + 1) % sieves.size
		} finally {
			lock.unlock()
		}
	}
	
	constructor(sieves: List<Sieve>) {
		this.sieves = Array<Sieve?>(sieves.size, sieves::get).copyOf(sieves.size + 1)
		obtainIdx = 0
		releaseIdx = sieves.size
		lock = ReentrantLock()
	}
}
