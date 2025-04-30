# 🏎️ F1 StatsHub System

F1 StatsHub es una API REST desarrollada en Go que expone estadísticas de la temporada 2024 de Fórmula 1 a través de endpoints organizados por corredores, carreras y resumen general.

## 📁 Estructura del Proyecto

```
F1_StatsHub/
├── db/
│   ├── loader.go
│   └── tablas.go
├── endpoints/
│   ├── carrera_detalle.go
│   ├── carreras.go
│   ├── drivers.go
│   ├── resumen.go
│   └── router.go
├── models/
│   └── models.go
├── proxy.db
├── server.go
├── main.go
├── go.mod / go.sum
```

## 🚀 Cómo ejecutar el proyecto

### 1. Requisitos

- Go 1.20+
- SQLite
- Módulo `github.com/mattn/go-sqlite3`

### 2. Instrucciones

#### Modo 1: Con base de datos ya cargada

> ⚠️ Este es el modo recomendado para evaluación.

```bash
go run server.go
```

Esto levanta el servidor en:

```
http://localhost:8080
```

## 📡 Endpoints disponibles

| Endpoint                          | Descripción                                     |
|----------------------------------|-------------------------------------------------|
| `/api/corredor`                  | Lista todos los pilotos                        |
| `/api/corredor/detalle/:id`      | Muestra detalles y estadísticas de un piloto   |
| `/api/carrera`                   | Muestra información básica de cada carrera     |
| `/api/carrera/detalle/:id`       | Entrega detalles como podios y vuelta rápida   |
| `/api/resumen`                   | Muestra top 3 de ganadores, pole y fastest lap |

## 🧠 Consideraciones

- **No es necesario conectarse a Internet**, ya que los datos están precargados en `proxy.db`.
- Todos los endpoints utilizan consultas SQL optimizadas.
- Los handlers están organizados modularmente por archivo para claridad y escalabilidad.

## 🙋 Autor

Carlos Ramírez  
Ingeniería Civil Informática - UTFSM  
Curso Sistemas Distribuidos
