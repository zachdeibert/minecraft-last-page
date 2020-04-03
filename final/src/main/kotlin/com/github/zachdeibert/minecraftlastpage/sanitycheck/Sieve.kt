package com.github.zachdeibert.minecraftlastpage.sanitycheck

import com.github.zachdeibert.minecraftlastpage.core.DimHash
import com.github.zachdeibert.minecraftlastpage.core.LastPage

class Sieve {
	companion object {
		val ALPHABET: CharArray = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()".toCharArray()
		
		@Synchronized
		private fun printf(format: String, vararg args: Any) {
			System.out.printf(format, *args)
		}
	}
	
	private var prefix: String
	private val max: Int
	private var str: CharArray
	private var idxs: IntArray
	var count: Long
	
	override fun toString(): String {
		return String.format("Sieve [prefix='%s', max=%d]: {str='%s'}", prefix, max, String(str))
	}
	
	fun run(_iters: Int): Boolean {
		var iters = _iters
		search@ while (str.size <= max && --iters >= 0) {
			for (i in idxs.size - 1 downTo 0) {
				if (++idxs[i] < ALPHABET.size) {
					str[prefix.length + i] = ALPHABET[idxs[i]]
					val s = String(str)
					++count
					if (DimHash.getHash(s) == LastPage.DIM_HASH) {
						printf("Matched hash: '%s'\n", s)
						val res = LastPage.decodeMessage(s)
						if (res != null) {
							val nice = res.joinToString("\n")
							printf("Encryption hacked: '%s'\nMessage:\n%s\n", s, nice)
						}
					}
					continue@search
				} else {
					idxs[i] = 0
					str[prefix.length + i] = ALPHABET[0]
				}
			}
			str = str.copyOf(str.size + 1)
			idxs = idxs.copyOf(idxs.size + 1)
			idxs[idxs.size - 1] = -1
		}
		return str.size > max
	}
	
	constructor(prefix: String, max: Int) {
		this.prefix = prefix;
		this.max = max;
		str = CharArray(prefix.length + 1)
		idxs = intArrayOf(-1)
		prefix.toCharArray().copyInto(str)
		count = 0
	}
}
