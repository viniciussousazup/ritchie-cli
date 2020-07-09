package template

const (
	StartFile = "Main"

	Main = `package com.ritchie.formula;

import com.ritchie.formula.{{bin-name}}.{{bin-name-first-upper}};

public class Main {

    public static void main(String[] args) throws Exception {
        String input1 = System.getenv("SAMPLE_TEXT");
        String input2 = System.getenv("SAMPLE_LIST");
        boolean input3 = Boolean.parseBoolean(System.getenv("SAMPLE_BOOL"));
        {{bin-name-first-upper}} {{bin-name}} = new {{bin-name-first-upper}}(input1, input2, input3);
        {{bin-name}}.Run();
    }
}`

	Dockerfile = `
FROM maven:3.6.3-jdk-8 AS builder

ADD . /app
WORKDIR /app
RUN mvn clean install


FROM alpine:latest
USER root

COPY --from=builder /app/target/Main.jar Main.jar
COPY --from=builder /app/set_umask.sh set_umask.sh

RUN apk update
RUN apk fetch openjdk8
RUN apk add openjdk8

ENV JAVA_HOME=/usr/lib/jvm/java-1.8-openjdk
ENV PATH="$JAVA_HOME/bin:${PATH}"

RUN chmod +x set_umask.sh

WORKDIR /app

ENTRYPOINT ["../set_umask.sh"]

CMD ["java -jar ../Main.jar"]`

	File = `package com.ritchie.formula.{{bin-name}};

public class {{bin-name-first-upper}} {

    private String input1;
    private String input2;
    private boolean input3;

    public void Run() throws Exception {
        System.out.printf("Hello World!\n");
        System.out.printf("You receive %s in text.\n", input1);
        System.out.printf("You receive %s in list.\n", input2);
        System.out.printf("You receive %s in boolean.\n", input3);
    }

    public {{bin-name-first-upper}}(String input1, String input2, boolean input3) {
        this.input1 = input1;
        this.input2 = input2;
        this.input3 = input3;
    }

    public String getInput1() {
        return input1;
    }

    public void setInput1(String input1) {
        this.input1 = input1;
    }

    public String getInput2() {
        return input2;
    }

    public void setInput2(String input2) {
        this.input2 = input2;
    }

    public boolean isInput3() {
        return input3;
    }

    public void setInput3(boolean input3) {
        this.input3 = input3;
    }
}`

	Makefile = `# Build parameters
BIN_FOLDER=../bin
SH=$(BIN_FOLDER)/run.sh
BAT=$(BIN_FOLDER)/run.bat
JAR_NAME=Main.jar

build: mvn-build sh-unix bat-windows

mvn-build:
	mkdir -p $(BIN_FOLDER)
	mvn clean install
	cp target/$(JAR_NAME) $(BIN_FOLDER)/$(JAR_NAME)
	#Clean files
	rm -Rf target

sh-unix:
	echo '#!/bin/sh' > $(SH)
	echo 'java -jar $(JAR_NAME)' >> $(SH)
	chmod +x $(SH)

bat-windows:
	echo '@ECHO OFF' > $(BAT)
	echo 'java -jar $(JAR_NAME)' >> $(BAT)`

	WindowsBuild = `:: Java parameters
echo off
SETLOCAL
SET BIN_FOLDER=..\bin
SET BIN_NAME=Main.jar
SET BAT_FILE=%BIN_FOLDER%\run.bat

:build
    call mvn clean install
    mkdir %BIN_FOLDER%
    copy target\%BIN_NAME% %BIN_FOLDER%\%BIN_NAME%
    del /Q /F target
    GOTO BAT_WINDOWS
    GOTO DONE

:BAT_WINDOWS
    	echo @ECHO OFF > %BAT_FILE%
    	echo java -jar %BIN_NAME% >> %BAT_FILE%
:DONE`

	Pom = `
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>com.ritchie.formula</groupId>
    <artifactId>#rit{{artifactId}}</artifactId>
    <version>1.0-SNAPSHOT</version>

    <properties>
        <maven.compiler.source>1.8</maven.compiler.source>
        <maven.compiler.target>1.8</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <maven-jar-plugin.version>3.2.0</maven-jar-plugin.version>
    </properties>

    <build>
        <finalName>Main</finalName>
        <plugins>
            <plugin>
                <!-- Build an executable JAR -->
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-jar-plugin</artifactId>
                <version>${maven-jar-plugin.version}</version>
                <configuration>
                    <archive>
                        <manifest>
                            <!-- <addClasspath>true</addClasspath> -->
                            <mainClass>com.ritchie.formula.Main</mainClass>
                        </manifest>
                    </archive>
                </configuration>
            </plugin>
        </plugins>
    </build>

    <dependencies>
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>4.12</version>
            <scope>test</scope>
        </dependency>
    </dependencies>
</project>
`
)
