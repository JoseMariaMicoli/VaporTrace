import markdown2
import pdfkit
import sys
import os
import re
import shutil

def patch_content_for_report(md_content):
    """
    Inyecta el Executive Summary en inglés y limpia glifos incompatibles.
    """
    # 1. Definición del Executive Summary basado en los hallazgos del PDF (CVE-2022-23131, etc.)
    executive_summary = """
## 1. EXECUTIVE SUMMARY

### 1.1 Risk Overview
VaporTrace Tactical Suite performed an automated adversarial emulation against the target infrastructure. The assessment identified critical bypass vectors and infrastructure-level weaknesses.

**OVERALL RISK RATING:** <span class="critical-text">CRITICAL</span>

| METRIC | VALUE |
| :--- | :--- |
| **Total Findings** | 59 |
| **Unique Targets** | 21 |
| **Average CVSS** | 1.4 / 10.0 |

**Key Vulnerability Insights:**
* **Authentication Bypass (CVE-2022-23131):** Critical session manipulation vulnerability found in administrative endpoints.
* **Infrastructure DoS (CVE-2023-44487):** High risk of service disruption via HTTP/2 Rapid Reset.
* **API Security:** Systematic failures in Property Level Authorization (API3:2023).

---
"""
    # Reemplazamos la sección 1 antigua por la nueva en inglés
    md_content = re.sub(r'#.*EXECUTIVE SUMMARY.*?(?=## 2)', executive_summary, md_content, flags=re.DOTALL)

    # 2. Reemplazo de barras de progreso por HTML con color para el PDF
    md_content = re.sub(r'\(█+░+\)', r'<span class="progress-bar">█</span>', md_content)
    md_content = re.sub(r'\(░+\)', r'<span class="progress-bar empty">█</span>', md_content)
    
    return md_content

def convert_md_to_pdf(input_md, output_pdf):
    try:
        with open(input_md, 'r', encoding='utf-8') as f:
            md_content = f.read()
    except FileNotFoundError:
        print(f"❌ Error: El archivo '{input_md}' no existe.")
        return

    # Aplicamos los parches de contenido
    md_content = patch_content_for_report(md_content)

    # Convertir MD a HTML
    html_content = markdown2.markdown(md_content, extras=[
        "tables", "fenced-code-blocks", "break-on-newline", "cuddled-lists"
    ])

    # CSS mejorado para visualización profesional y colores de severidad
    css = """
    <style>
        @page { margin: 2cm; }
        body { font-family: 'Helvetica', sans-serif; line-height: 1.6; color: #1a1a1a; font-size: 11pt; }
        h1 { color: #d73a49; border-bottom: 2px solid #d73a49; padding-bottom: 10px; }
        h2 { color: #24292e; border-bottom: 1px solid #eaecef; padding-bottom: 5px; margin-top: 30px; }
        
        table { border-collapse: collapse; width: 100%; margin: 20px 0; }
        th, td { border: 1px solid #dfe2e5; padding: 10px; text-align: left; }
        th { background-color: #f6f8fa; font-weight: bold; }
        
        .critical-text { color: #cf222e; font-weight: bold; }
        .progress-bar { color: #cf222e; font-family: monospace; font-weight: bold; }
        .progress-bar.empty { color: #dfe2e5; }
        
        pre { background-color: #f6f8fa; padding: 15px; border-radius: 5px; border: 1px solid #ddd; }
        code { background-color: #f3f3f3; padding: 2px 4px; border-radius: 3px; font-family: monospace; }
    </style>
    """
    
    full_html = f"<html><head><meta charset='UTF-8'>{css}</head><body>{html_content}</body></html>"

    # Configuración de binarios
    wk_path = shutil.which("wkhtmltopdf") or '/usr/bin/wkhtmltopdf'
    config = pdfkit.configuration(wkhtmltopdf=wk_path)

    options = {
        'encoding': "UTF-8",
        'quiet': '',
        'enable-local-file-access': '',
        'margin-top': '0.75in',
        'margin-right': '0.75in',
        'margin-bottom': '0.75in',
        'margin-left': '0.75in',
    }

    try:
        pdfkit.from_string(full_html, output_pdf, configuration=config, options=options)
        print(f"✅ Reporte Final Generado (Executive Summary en Inglés): {output_pdf}")
    except Exception as e:
        print(f"❌ Error: {e}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Uso: python md2pdf.py <archivo.md>")
        sys.exit(1)

    in_f = sys.argv[1]
    out_f = os.path.splitext(in_f)[0] + ".pdf"
    convert_md_to_pdf(in_f, out_f)