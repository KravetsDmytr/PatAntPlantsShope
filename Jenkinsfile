pipeline {
    agent any
    environment {
        WEB_PORT = '18080'
        POSTGRES_PORT = '55432'
        API_BASE_URL = 'http://host.docker.internal:18080/api/v1'
    }

    stages {
        stage('Start API for tests') {
            steps {
                bat 'docker compose up --build -d postgres web'
            }
        }

        stage('Run tests in Docker Python') {
            steps {
                bat 'docker run --rm -e API_BASE_URL=%API_BASE_URL% -v "%CD%:/app" -w /app python:3.12-slim sh -c "pip install -r requirements.txt && pytest tests --junitxml=report.xml -v"'
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