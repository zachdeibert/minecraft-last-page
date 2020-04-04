package com.github.zachdeibert.minecraftlastpage.checker

import java.io.Closeable
import java.util.Scanner
import java.io.FileInputStream
import java.util.stream.Stream
import java.util.NoSuchElementException

class InputFile : Closeable, Iterator<String>, Iterable<String> {
	private val scan: Scanner
	private var next: String?
	
	override fun hasNext(): Boolean {
		if (next == null) {
			try {
				next = scan.nextLine()
			} catch (ex: Exception) {
				return false
			}
		}
		return true
	}

	override fun next(): String {
		if (hasNext()) {
			val s = next
			next = null
			return s!!
		} else {
			throw NoSuchElementException()
		}
	}
	
	override fun iterator(): Iterator<String> {
		return this
	}
	
	override fun close() {
		scan.close()
	}
	
	constructor(filename: String) {
		scan = Scanner(FileInputStream(filename))
		next = null
	}
}
