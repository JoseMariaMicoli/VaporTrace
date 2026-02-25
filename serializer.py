"""
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
"""
/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

import os

# Configuraciones de filtrado
ALLOWED_EXTENSIONS = {'.go', '.mod', '.sum'}
EXCLUDED_DIRS = {
    '.git', 'vendor', 'bin', 'pkg/mod', 'obj', 
    '.vscode', '.idea', '__pycache__', 'testdata'
}

def is_binary(file_path):
    """Comprueba si un archivo es binario leyendo los primeros bytes."""
    try:
        with open(file_path, 'rb') as f:
            chunk = f.read(1024)
            return b'\x00' in chunk # Si tiene bytes nulos, es binario
    except:
        return True

def serialize_project():
    output_file = "codebase_vaportrace.txt"
    
    # Verificar si estamos en la raíz de un proyecto Go
    if not os.path.exists("go.mod"):
        print("[ERROR] No se encontró 'go.mod'. Ejecuta el script en la raíz de VaporTrace.")
        return

    print(f"--- Iniciando Serialización de VaporTrace ---")
    
    with open(output_file, 'w', encoding='utf-8') as out:
        file_count = 0
        for root, dirs, files in os.walk('.'):
            # Filtrar directorios prohibidos
            dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS]
            
            for file in files:
                ext = os.path.splitext(file)[1].lower()
                file_path = os.path.join(root, file)
                
                # Normalizar ruta para mostrarla en el archivo
                rel_path = os.path.relpath(file_path, '.')

                if ext in ALLOWED_EXTENSIONS:
                    if not is_binary(file_path):
                        try:
                            with open(file_path, 'r', encoding='utf-8') as f:
                                content = f.read()
                                out.write(f"\n{'#'*60}\n")
                                out.write(f"PATH: {rel_path}\n")
                                out.write(f"{'#'*60}\n\n")
                                out.write(content)
                                out.write("\n")
                                file_count += 1
                                print(f"[+] Incluido: {rel_path}")
                        except Exception as e:
                            print(f"[!] Error leyendo {rel_path}: {e}")

    print(f"\n[ÉXITO] {file_count} archivos serializados en: {output_file}")
    print("Sube este archivo a Gemini 1.5 Pro junto con el prompt de arquitectura.")

if __name__ == "__main__":
    serialize_project()