package com.github.zachdeibert.minecraftlastpage.sanitycheck

class Worker : Thread {
	private val pool: Pool
	private val time: Int
	
	override fun run() {
		try {
			while (true) {
				val sieve: Sieve? = pool.obtain()
				if (sieve == null) {
					break
				} else if (!sieve.run(time)) {
					pool.release(sieve)
				}
			}
			System.out.println("Killing worker thread due to no work to do")
		} catch (t: Throwable) {
			t.printStackTrace()
		}
	}
	
	constructor(pool: Pool, time: Int) {
		this.pool = pool
		this.time = time
	}
}
