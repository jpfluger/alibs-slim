# Markdown Testing for Goldmark

This document contains various Markdown elements for testing with the Goldmark Markdown-to-HTML converter.

---

## Headers

### H1 Header
# This is an H1 Header

### H2 Header
## This is an H2 Header

### H3 Header
### This is an H3 Header

---

## Emphasis

*Italic Text*  
**Bold Text**  
***Bold and Italic Text***  

---

## Lists

### Unordered List
- Item 1
- Item 2
  - Subitem 2.1
  - Subitem 2.2

### Ordered List
1. First item
2. Second item
   1. Subitem 2.1
   2. Subitem 2.2

---

## Links and Images

### Links
[OpenAI](https://www.openai.com)  
[Local File](./b.js)  
[HTML Embedded File](./c.html)

### Images
![Test Image](https://via.placeholder.com/150)

---

## Tables

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Row 1 Col 1 | Row 1 Col 2 | Row 1 Col 3 |
| Row 2 Col 1 | Row 2 Col 2 | Row 2 Col 3 |

---

## Code Blocks

### Inline Code
This is `inline code`.

### Code Block with Backticks
```javascript
console.log("This is a JavaScript code block.");
```

```python
print("This is a Python code block.")
```

---

## Blockquotes

> This is a blockquote.  
> It can span multiple lines.

---

## Embedded JavaScript and HTML

### Linked JavaScript
<script src="./b.js"></script>

### Embedded JavaScript
```html
<script>
  console.log("This is an embedded JavaScript block.");
</script>
```

### Linked HTML
<a href="./c.html">Click here to see the embedded HTML file.</a>

### Embedded HTML
```html
<div>
  <h2>Embedded HTML Example</h2>
  <p>This is an HTML snippet embedded directly into Markdown.</p>
</div>
```

---

## Task Lists

- [x] Completed Task
- [ ] Incomplete Task
- [ ] Another Incomplete Task

---

## Horizontal Rules

---

---

---

## Miscellaneous

### HTML Entities
&copy; &reg; &amp; &lt; &gt;

### Footnotes
This is a sentence with a footnote.[^1]

[^1]: This is the footnote text.
