# JotterXpress - RPM Package

Esta guía explica cómo crear e instalar un paquete RPM de JotterXpress.

## Requisitos Previos

### Instalar herramientas de desarrollo
```bash
# En RHEL/CentOS/Fedora
sudo dnf install -y rpm-build rpmdevtools golang git
# o en versiones anteriores:
sudo yum install -y rpm-build rpmdevtools golang git
```

### Configurar el directorio de build (opcional pero recomendado)
```bash
# Esto crea ~/rpmbuild con la estructura correcta
rpmdev-setuptree
```

## Construcción del RPM

### Método 1: Usando el script automático (Recomendado)

```bash
# Ejecuta el script de empaquetado
./package-rpm.sh
```

Este script automáticamente:
1. Compila la aplicación
2. Crea un tarball de los fuentes
3. Configura el directorio de build RPM
4. Construye el paquete RPM

### Método 2: Manual

Si prefieres hacerlo manualmente:

```bash
# 1. Compilar la aplicación
go build -o jotterxpress cmd/jotterxpress/main.go

# 2. Crear tarball de fuentes
tar --exclude='.git' \
    --exclude='rpmbuild' \
    --exclude='*.tar.gz' \
    --exclude='jotterxpress.spec' \
    --exclude='package-rpm.sh' \
    --exclude='bin' \
    -czf jotterxpress-1.0.0.tar.gz .

# 3. Copiar archivos a ~/rpmbuild/
cp jotterxpress-1.0.0.tar.gz ~/rpmbuild/SOURCES/
cp jotterxpress.spec ~/rpmbuild/SPECS/

# 4. Construir el RPM
cd ~/rpmbuild/SPECS
rpmbuild -ba jotterxpress.spec
```

## Instalación

### Instalar el RPM
```bash
# Instalar desde el archivo generado
sudo rpm -ivh ~/rpmbuild/RPMS/*/jotterxpress-1.0.0-*.rpm

# O si ya existe una versión anterior:
sudo rpm -Uvh ~/rpmbuild/RPMS/*/jotterxpress-1.0.0-*.rpm
```

### Verificar instalación
```bash
# Verificar que el comando está disponible
jtx --version

# Ver archivos instalados
rpm -ql jotterxpress

# Ver información del paquete
rpm -qi jotterxpress
```

## Uso del paquete instalado

Una vez instalado, puedes usar JotterXpress:

```bash
# Crear una nota rápida
jtx "Mi primera nota"

# Listar notas de hoy
jtx --list

# Crear una tarea
jtx --task

# Ver ayuda
jtx --help
```

## Desinstalación

```bash
sudo rpm -e jotterxpress
```

## Personalización del Spec File

Puedes personalizar `jotterxpress.spec` para:
- Cambiar la versión (variable `Version`)
- Cambiar la descripción (`%description`)
- Agregar dependencias adicionales (`BuildRequires`, `Requires`)
- Modificar rutas de instalación
- Agregar scripts pre/post-installación

## Solución de Problemas

### Error: "command 'go' not found"
```bash
# Asegúrate de que Go esté instalado y en tu PATH
which go
go version
```

### Error: "rpmbuild: command not found"
```bash
# Instala rpm-build y rpmdevtools
sudo dnf install rpm-build rpmdevtools
```

### Error: "directory does not exist: BUILD/BUILDROOT/..."
```bash
# Crea la estructura de directorios
mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
# O usa:
rpmdev-setuptree
```

### El RPM se construye pero no instala
```bash
# Verifica las dependencias faltantes
sudo rpm -ivh --test ~/rpmbuild/RPMS/*/jotterxpress-*.rpm

# Verifica errores específicos
sudo rpm -ivh -vv ~/rpmbuild/RPMS/*/jotterxpress-*.rpm
```

## Archivos Generados

Después de la construcción exitosa, encontrarás:

- **RPM binario**: `~/rpmbuild/RPMS/*/jotterxpress-1.0.0-*.rpm`
- **RPM fuente**: `~/rpmbuild/SRPMS/jotterxpress-1.0.0-*.src.rpm`
- **Tarball fuente**: `jotterxpress-1.0.0.tar.gz`

## Compartir el RPM

Puedes compartir el archivo `.rpm` con otros usuarios. Ellos pueden instalarlo con:

```bash
sudo rpm -ivh jotterxpress-1.0.0-*.rpm
```

## Distribuir a través de un repositorio

Si quieres crear un repositorio RPM personalizado:

```bash
# Crear repositorio
mkdir -p /var/www/html/repos/jotterxpress
cp ~/rpmbuild/RPMS/*/jotterxpress-*.rpm /var/www/html/repos/jotterxpress/

# Crear metadatos del repositorio
createrepo /var/www/html/repos/jotterxpress/

# Los usuarios pueden agregar esto a su /etc/yum.repos.d/
# [jotterxpress]
# name=JotterXpress Repository
# baseurl=http://your-server/repos/jotterxpress/
# enabled=1
# gpgcheck=0
```

## Soporte

Para problemas o preguntas, abre un issue en el repositorio del proyecto.
