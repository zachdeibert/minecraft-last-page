package com.github.zachdeibert.minecraftlastpage.checker

import com.github.zachdeibert.minecraftlastpage.core.DimHash
import com.github.zachdeibert.minecraftlastpage.core.LastPage

fun main(args: Array<String>) {
	val input = InputFile("../candidates.txt")
	input.use {
		input.forEach {
			if (DimHash.getHash(it) != LastPage.DIM_HASH) {
				System.out.printf("WARNING: hash didn't match for entry '%s'\n", it)
			} else {
				val res = LastPage.decodeMessage(it)
				if (res != null) {
					val nice = res.joinToString("\n")
					System.out.printf("Encryption hacked: '%s'\nMessage:\n%s\n", it, nice)
				}
			}
		}
	}
}
