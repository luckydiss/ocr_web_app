const MarkdownRenderer = {
    init() {
        if (typeof marked !== 'undefined') {
            marked.setOptions({
                breaks: true,
                gfm: true
            });
        }
    },

    renderLatex(text, element) {
        if (typeof katex === 'undefined') {
            element.textContent = text;
            return;
        }

        let html = text;

        html = html.replace(/\$\$([\s\S]+?)\$\$/g, (match, tex) => {
            try {
                return katex.renderToString(tex.trim(), {
                    displayMode: true,
                    throwOnError: false,
                    output: 'html'
                });
            } catch (e) {
                console.error('KaTeX display error:', e);
                return match;
            }
        });

        html = html.replace(/\$([^\$\n]+?)\$/g, (match, tex) => {
            try {
                return katex.renderToString(tex.trim(), {
                    displayMode: false,
                    throwOnError: false,
                    output: 'html'
                });
            } catch (e) {
                console.error('KaTeX inline error:', e);
                return match;
            }
        });

        element.innerHTML = html;
    },

    render(markdown, container) {
        if (!markdown) {
            container.innerHTML = '<p style="color: var(--hint-color)">No content extracted</p>';
            return;
        }

        const tempDiv = document.createElement('div');

        const mathPlaceholders = [];
        let processedMarkdown = markdown;

        processedMarkdown = processedMarkdown.replace(/\$\$([\s\S]+?)\$\$/g, (match, tex) => {
            const placeholder = `%%DISPLAY_MATH_${mathPlaceholders.length}%%`;
            mathPlaceholders.push({ type: 'display', tex: tex.trim() });
            return placeholder;
        });

        processedMarkdown = processedMarkdown.replace(/\$([^\$\n]+?)\$/g, (match, tex) => {
            const placeholder = `%%INLINE_MATH_${mathPlaceholders.length}%%`;
            mathPlaceholders.push({ type: 'inline', tex: tex.trim() });
            return placeholder;
        });

        if (typeof marked !== 'undefined') {
            tempDiv.innerHTML = marked.parse(processedMarkdown);
        } else {
            tempDiv.innerHTML = processedMarkdown.replace(/\n/g, '<br>');
        }

        let html = tempDiv.innerHTML;

        mathPlaceholders.forEach((item, index) => {
            const displayPlaceholder = `%%DISPLAY_MATH_${index}%%`;
            const inlinePlaceholder = `%%INLINE_MATH_${index}%%`;

            if (typeof katex !== 'undefined') {
                try {
                    const rendered = katex.renderToString(item.tex, {
                        displayMode: item.type === 'display',
                        throwOnError: false,
                        output: 'html'
                    });
                    html = html.replace(displayPlaceholder, rendered);
                    html = html.replace(inlinePlaceholder, rendered);
                } catch (e) {
                    const original = item.type === 'display' ? `$$${item.tex}$$` : `$${item.tex}$`;
                    html = html.replace(displayPlaceholder, original);
                    html = html.replace(inlinePlaceholder, original);
                }
            } else {
                const original = item.type === 'display' ? `$$${item.tex}$$` : `$${item.tex}$`;
                html = html.replace(displayPlaceholder, original);
                html = html.replace(inlinePlaceholder, original);
            }
        });

        container.innerHTML = html;
    }
};

MarkdownRenderer.init();
