plugins {
    java
}

group = "com.wbrawner"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

java.sourceCompatibility = JavaVersion.VERSION_11
java.targetCompatibility = JavaVersion.VERSION_11

val invoker by configurations.creating

dependencies {
    implementation("com.google.cloud:google-cloud-dns:2.0.3")
    invoker("com.google.cloud.functions.invoker:java-function-invoker:1.1.0")
    compileOnly("com.google.cloud.functions:functions-framework-api:1.0.4")
    testImplementation("com.google.cloud.functions:functions-framework-api:1.0.4")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.8.2")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine")
}

tasks.getByName<Test>("test") {
    useJUnitPlatform()
}

tasks.register<JavaExec>("runFunction") {
    mainClass.set("com.google.cloud.functions.invoker.runner.Invoker")
    classpath(configurations.getByName("invoker"))
    inputs.files(configurations.runtimeClasspath, sourceSets.main.get().output)
    args(
            "--target", project.findProperty("run.functionTarget") ?: "",
            "--port", project.findProperty("run.port") ?: 8080
    )
    doFirst {
        args("--classpath", files(configurations.runtimeClasspath, sourceSets.main.get().output).asPath)
    }
}
