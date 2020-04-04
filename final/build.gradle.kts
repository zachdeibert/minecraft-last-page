plugins {
    base
    eclipse
    kotlin("jvm") version "1.3.71"
}

repositories {
    mavenCentral()
}

dependencies {
    implementation(kotlin("stdlib-jdk8"))
    implementation("com.google.guava:guava:28.0-jre")
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

tasks {
    register<JavaExec>("runSanityCheck") {
        main = "com.github.zachdeibert.minecraftlastpage.sanitycheck.SanityCheckKt"
        classpath = sourceSets["main"].runtimeClasspath
    }

    register<JavaExec>("runChecker") {
        main = "com.github.zachdeibert.minecraftlastpage.checker.CheckerKt"
        classpath = sourceSets["main"].runtimeClasspath
    }
}
