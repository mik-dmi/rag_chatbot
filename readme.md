# Retrieval-Augmented Generation (RAG) in a RESTful API  
**Technologies:** Golang, Docker, Weaviate, Redis, Langchain

## What is RAG?  
In simple terms, Retrieval-Augmented Generation (RAG) is an AI system that improves chatbot responses by first searching for relevant information from external data sources and then using that data to generate more accurate and reliable answers.

## Project Goal  
The goal of this project is to create a tool that simplifies working with large documentation. By using a RAG system, users can easily navigate extensive documents. The aim is to develop an AI chatbot that not only answers users’ questions related to specific documents but also directs them to the exact location in the documentation where the information can be found.

## Features  
This RAG system is implemented as a RESTful API, allowing users to update the vector database as well as interact directly with the RAG system. The application uses Weaviate as a vector database and Redis for storage — both of which can be run locally. The API includes standard functionality and security measures. Some of the features implemented (or planned for the near future) include:
- **Rate limiting**
- **Locally run databases**
- **Repository pattern** for flexibility and the ability to easily swap out databases
- **Reliable logging**
- **Docker support** for running on your machine
- **Chatbot memory**
- **Easy configuration** for changing the underlying LLM (Language Model)
- **Readable and clean code**
- **Open source**
- **Standardized API error responses**
- **Authentication**

## Project Overview  
Below is a diagram providing a simplified overview of the entire API and RAG system, which helps explain how everything works.

![RAG API Diagram](./public/rag_diagram.png)
