# 📝 JotterXpress

Una herramienta CLI rápida y simple para tomar notas, construida en Go con arquitectura hexagonal.

## 🚀 Características

- **Toma de notas rápida**: Crea notas con un simple comando
- **Tipos de notas**: Soporta notas de texto, tareas y contactos
- **Interfaz interactiva**: Lista navegable con Bubble Tea
- **Organización por fechas**: Cada día se genera un archivo JSON separado
- **Listado de notas**: Visualiza todas tus notas del día
- **Prioridades de tareas**: Sistema de prioridades (low, medium, high, urgent)
- **Gestión de contactos**: Almacena nombres, teléfonos y emails
- **Almacenamiento JSON**: Formato estructurado y extensible
- **Detección automática**: Interfaz interactiva o texto según el entorno
- **Arquitectura hexagonal**: Código limpio y mantenible
- **Interfaz hermosa**: Colores y emojis para una mejor experiencia

## 📦 Instalación

### Opción 1: Compilar desde el código fuente

```bash
# Clonar el repositorio
git clone <tu-repo>
cd jotterxpress

# Compilar
make build

# Instalar (opcional)
make install
```

### Opción 2: Usar directamente

```bash
# Compilar y ejecutar
make quick ARGS="\"Hola, esta es mi primera nota\""
```

## 🎯 Uso

### Crear una nota de texto
```bash
# Modo normal (línea de comandos)
jtx "Esta es mi nota de hoy"

# Modo interactivo (textarea)
jtx -i
```

### Crear una tarea
```bash
# Modo interactivo (por defecto)
jtx task

# Modo online (línea de comandos)
jtx task "Completar proyecto" --priority high --online
```

### Crear un contacto
```bash
# Modo interactivo (por defecto)
jtx contact

# Modo online (línea de comandos)
jtx contact "Juan Pérez" --phone "+1234567890" --email "juan@example.com" --online
```

### Listar notas del día (interactivo)
```bash
jtx list
```

### Listar notas de una fecha específica (interactivo)
```bash
jtx list-date 2024-01-15
```

### Modo interactivo dedicado
```bash
jtx interactive
```

### Ver ayuda
```bash
jtx --help
```

### Ver ayuda de comandos específicos
```bash
jtx task --help
jtx contact --help
```

## 📋 Formularios Interactivos (Modo por Defecto)

### Formulario de Tareas
- **Campos**: Descripción, prioridad, asignado, horas estimadas
- **Validación**: Prioridades válidas, números para horas
- **Navegación**: Tab/Shift+Tab para moverse entre campos
- **Uso**: `jtx task` (modo interactivo por defecto)

### Formulario de Contactos
- **Campos**: Nombre, teléfono, email, dirección
- **Validación**: Formato de teléfono y email
- **Requerimientos**: Al menos teléfono o email
- **Uso**: `jtx contact` (modo interactivo por defecto)

### Controles del Formulario
- `Tab` o `Ctrl+N`: Siguiente campo
- `Shift+Tab` o `Ctrl+P`: Campo anterior
- `Enter`: Completar formulario
- `Ctrl+C` o `Esc`: Cancelar

### Modo Online (Línea de Comandos)
- **Tareas**: `jtx task "descripción" --priority high --online`
- **Contactos**: `jtx contact "nombre" --phone "+1234567890" --online`

### Textarea Interactivo para Notas
- **Comando**: `jtx -i`
- **Características**: Textarea de 80x15 caracteres, límite de 2000 caracteres
- **Controles**: `Ctrl+J` para guardar, `Ctrl+C` para cancelar
- **Uso**: Ideal para notas largas con múltiples líneas

## 🎮 Interfaz Interactiva

### Controles de Navegación
- `↑/↓` o `j/k`: Navegar por la lista
- `Space`: Seleccionar/deseleccionar nota
- `Enter`: Seleccionar nota
- `q` o `Ctrl+C`: Salir
- `Esc`: Salir

### Características Visuales
- **Colores por tipo**: Tareas (verde), contactos (morado), texto (blanco)
- **Navegación fluida**: Cursor visual y selección múltiple
- **Detección automática**: Interfaz interactiva en TTY, texto en pipelines
- **Fallback inteligente**: Si falla la interfaz interactiva, usa modo texto

## 📁 Estructura del Proyecto

```
jotterxpress/
├── cmd/jotterxpress/          # Punto de entrada de la aplicación
├── internal/
│   ├── domain/               # Lógica de dominio (entidades y puertos)
│   │   ├── entities/         # Entidades del dominio
│   │   └── ports/           # Interfaces (contratos)
│   ├── application/          # Lógica de aplicación
│   │   └── services/        # Servicios de aplicación
│   └── adapters/            # Adaptadores (implementaciones)
│       ├── repository/      # Adaptador de repositorio (archivos)
│       └── cli/            # Adaptador CLI
├── go.mod                   # Dependencias de Go
├── Makefile                # Comandos de build
└── README.md              # Este archivo
```

## 🏗️ Arquitectura

Este proyecto utiliza **Arquitectura Hexagonal** (Ports and Adapters), que proporciona:

- **Separación clara** entre la lógica de negocio y la infraestructura
- **Testabilidad** mejorada
- **Flexibilidad** para cambiar implementaciones
- **Mantenibilidad** del código

### Componentes principales:

1. **Domain Layer**: Contiene las entidades (`Note`) y puertos (interfaces)
2. **Application Layer**: Contiene la lógica de aplicación (servicios)
3. **Adapters Layer**: Implementaciones concretas (repositorio de archivos, CLI)

## 📝 Almacenamiento

Las notas se guardan en archivos JSON en el directorio:
```
~/.jotterxpress/notes/
├── 2024-01-15.json
├── 2024-01-16.json
└── ...
```

Cada archivo contiene todas las notas del día correspondiente en formato JSON estructurado:
```json
[
  {
    "id": "1761259384208901547",
    "type": "text",
    "content": "Mi primera nota del día",
    "created_at": "2024-01-15T14:30:25Z",
    "updated_at": "2024-01-15T14:30:25Z",
    "date": "2024-01-15",
    "metadata": {}
  },
  {
    "id": "1761259388410638049",
    "type": "task",
    "content": "Completar proyecto",
    "created_at": "2024-01-15T15:45:12Z",
    "updated_at": "2024-01-15T15:45:12Z",
    "date": "2024-01-15",
    "metadata": {
      "priority": "high",
      "status": "pending"
    }
  }
]
```

## 🛠️ Desarrollo

### Comandos disponibles:

```bash
make build      # Compilar la aplicación
make install    # Instalar en /usr/local/bin
make clean      # Limpiar archivos de build
make test       # Ejecutar tests
make dev-setup  # Configurar entorno de desarrollo
make quick      # Compilar y ejecutar rápidamente
```

### Dependencias principales:

- **Cobra**: Framework CLI
- **Charm**: Librerías para interfaces hermosas
- **Viper**: Configuración (preparado para futuras extensiones)

## 🎨 Ejemplo de uso

```bash
# Crear una nota de texto
$ jx "Reunión importante a las 3 PM"
Note saved successfully!

# Crear una tarea
$ jx task "Completar proyecto JotterXpress" --priority high
Task created successfully!
📋 Type: task
📝 Content: Completar proyecto JotterXpress
⚡ Priority: high
📊 Status: pending

# Crear un contacto
$ jx contact "María García" --phone "+1234567890" --email "maria@example.com"
Contact created successfully!
👤 Type: contact
📝 Name: María García
📞 Phone: +1234567890
📧 Email: maria@example.com

# Listar notas del día
$ jx list
📝 Today's Notes

📝 Notes (3 found):

1. [14:30:25] Reunión importante a las 3 PM
2. [15:45:12] Completar proyecto JotterXpress [high, pending]
3. [16:20:30] María García [+1234567890]
```

## 🚀 Próximas características

- [ ] Búsqueda en notas
- [ ] Exportar notas
- [ ] Categorías/etiquetas
- [ ] Sincronización en la nube
- [ ] Editor interactivo
- [ ] Estadísticas de notas

## 📄 Licencia

MIT License - ver archivo LICENSE para más detalles.

---

**¡Disfruta tomando notas con JotterXpress! 🎉**
