from pdf2docx import Converter
import sys

def pdf_to_word(pdf_path, word_path):
    cv = Converter(pdf_path)
    cv.convert(word_path, start=0, end=None)
    cv.close()

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python pdf_to_word.py <input_pdf> <output_docx>")
        sys.exit(1)

    pdf_to_word(sys.argv[1], sys.argv[2])
