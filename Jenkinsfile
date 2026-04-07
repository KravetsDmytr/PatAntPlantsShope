pipeline {
    agent any

    stages {
        stage('Start API for tests') {
            steps {
                bat 'docker compose up --build -d postgres web'
            }
        }

        stage('Run tests in Docker Python') {
            steps {
                bat 'docker run --rm -v "%CD%:/app" -w /app python:3.12-slim sh -c "pip install -r requirements.txt && pytest tests --junitxml=report.xml -v"'
            }
        }
    }

    post {
        always {
            junit allowEmptyResults: true, testResults: 'report.xml'
            bat 'docker compose down'
        }
    }
}