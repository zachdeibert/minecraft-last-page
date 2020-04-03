package com.github.zachdeibert.minecraftlastpage.core

import com.google.common.hash.Hashing
import java.nio.charset.StandardCharsets

object DimHash {
	fun getHash(txt: String): Int {
		return Hashing.sha256().hashString(txt + ":why_so_salty#LazyCrypto", StandardCharsets.UTF_8).asInt() and Int.MAX_VALUE
	}
}
