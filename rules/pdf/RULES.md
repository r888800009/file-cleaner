# Rules

file_cleaner using yara rules to find the files that you want to move or delete.
for example we have a rule to match conference papers, and pdf can not be matched by text, so we have two rules to match pdf
- ~~`rules/pdf/papers.yar` only match pdf bytes, it will not match text in pdf.~~ (placeholder)
- `rules/pdf/papers_text.yar` match text in pdf, it would use `pdftotext` to extract text from pdf and match the text. maybe it is slow.
