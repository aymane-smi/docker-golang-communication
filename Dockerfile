FROM openjdk:11-jdk-slim

WORKDIR /app

CMD ["tail", "-f", "/dev/null"] 