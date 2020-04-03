package com.github.zachdeibert.minecraftlastpage.core

import javax.crypto.SecretKeyFactory
import java.security.spec.KeySpec
import javax.crypto.spec.PBEKeySpec
import javax.crypto.SecretKey
import javax.crypto.spec.SecretKeySpec
import javax.crypto.Cipher
import javax.crypto.spec.IvParameterSpec
import java.nio.ByteBuffer
import java.nio.charset.StandardCharsets

object LastPage {
	const val DIM_HASH: Int = 1791460938
	private val DATA: ByteArray = byteArrayOf(
			72, -71, 33, -116, 61, 25, -105, 61, 69, -34, 
			22, -96, 83, -41, 4, 100, 49, -120, -76, -32, 
			-24, 105, -103, 57, -101, -101, 114, 39, 45, -48, 
			-58, 106, -83, 72, -120, -98, 14, 111, -73, 38, 
			-43, -29, -17, -64, -48, -21, -63, -14, 7, -65, 
			-115, 61, -62, -121, 108, -2, 24, 84, -62, 117, 
			115, 52, -88, 18, 30, 115, 44, -113, -123, 88, 
			77, -122, 15, -13, 123, 85, -118, -12, -30, 7, 
			104, -75, -84, -57, -124, -113, -38, 84, 52, 6, 
			-94, 67, 76, 59, 105, 82, 92, -65, -52, -26, 
			-46, 45, 94, 47, 10, -14, -86, -12, 1, 111, 
			-107, -119, -115, 32, -60, -92, 107, -2, 73, 109, 
			Byte.MIN_VALUE, -107, -52, -7, -6, 126, 98, -71, -92, 12, 
			-41, 83, -124, 14, -51, -15, 4, -3, -65, -36, 
			99, 63, 119, 64, 46, 21, 10, 30, -23, -10, 
			-90, -36, -4, -106, -102, 84, -17, 58, 59, -76, 
			-103, -28, -95, 4, 112, 18, 3, -78, 125, -79, 
			11, 120, -59, -64, -37, -47, 19, -21, 90, -9, 
			-65, 109, 70, -83, -4, 34, 41, -109, 27, -20, 
			29, 60, 109, -117, 74, -112, -58, 76, 96, 9, 
			-65, 86, 63, 62, 112, -88, 96, -35, 64, 57, 
			35, 89, -24, -40, 121, 106, -102, -103, -24, -73, 
			103, -110, 56, 97, -82, 55, -53, -100, 22, -68, 
			104, 8, 98, -120, -65, -30, 38, 114, -59, 30, 
			66, -119, 59, -93, 107, -50, 115, 40, 80, 77, 
			-61, -102, -62, -110, -80, -85, 19, 123, -120, 70, 
			-119, 11, 63, 30, 92, 73, 81, -19, -14, 122, 
			-103, -108, 38, -116, -100, 50, -121, -7, -125, 61, 
			-44, -38, -117, 16, 14, -101, 79, -96, 89, 12, 
			84, -36, 42, -21, -109, -7, 117, 64, 38, 18, 
			-97, -58, 73, 2, 41, 70, -85, 75, 6, 123, 
			76, -66, 53, -41, 25, -14, -104, -19, 67, -28, 
			-9, -111, 59, -109, 35, 57, 108, 100, 40, 116, 
			-106, Byte.MIN_VALUE, 2, 109, -75, 3, 19, 87, -120, 59, 
			-20, -15, 74, -40, 106, -3, -122, 19, -94, 53, 
			-103, -60, -36, 2, 52, 31, 63, 17, -32, -61, 
			-116, 5, 9, 117, -72, -28, -125, 99, -54, -126, 
			96, 21, 29, 38, 35, 90, -32, 89, 48, 108, 
			10, -52, -117, 2, -74, -122, -21, 119, 126, -110, 
			-115, 57, -119, -53, 43, Byte.MIN_VALUE, 10, 97, 122, 126, 
			-111, 103, 113, 90, 101, 44, 9, 5, 102, 88, 
			-24, -108, -8, 42, 65, 46
	)
	private val IVS: ByteArray = byteArrayOf(
		-114, 123, -36, 36, 6, 2, 31, 116, -76, -125, 
		      -62, -61, -41, -121, 82, -106
	)
	
	fun decodeMessage(txt: String): List<String>? {
		try {
			val factory = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256")
			val spec = PBEKeySpec(txt.toCharArray(), "pinch_of_salt".toByteArray(StandardCharsets.UTF_8), 65536, 128)
			val key = factory.generateSecret(spec)
			val secret = SecretKeySpec(key.getEncoded(), "AES")
			val cipher = Cipher.getInstance("AES/CBC/PKCS5Padding")
			val iv = IvParameterSpec(IVS)
			cipher.init(Cipher.DECRYPT_MODE, secret, iv)
			val res = cipher.doFinal(DATA)
			return StandardCharsets.UTF_8.decode(ByteBuffer.wrap(res)).toString().split("\n")
		} catch (ex: Exception) {
			return null
		}
	}
}
