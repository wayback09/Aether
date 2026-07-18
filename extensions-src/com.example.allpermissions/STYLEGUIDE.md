# Style Guide for Extensions

This document outlines the style guide that should be followed when developing extensions for Aether. Adhering to these guidelines ensures consistency, maintainability, and a professional presentation of extensions.

## General Principles

1.  **Readability:** Code should be clean, well-formatted, and easy to understand. Use meaningful variable and function names.
2.  **Consistency:** Maintain a consistent coding style throughout your extension.
3.  **Modularity:** Break down complex logic into smaller, reusable functions or modules.
4.  **Error Handling:** Implement robust error handling to gracefully manage potential issues (e.g., network errors, file access problems).
5.  **User Experience:** For UI components, ensure they are intuitive, accessible, and visually consistent with the Aether launcher's design language.

## JavaScript (main.js and UI scripts)

*   **Indentation:** Use 2 spaces for indentation.
*   **Quotes:** Use single quotes (`'`) for strings unless the string itself contains a single quote, in which case use double quotes (`"`).
*   **Variable Declaration:** Use `const` by default, and `let` only when a variable needs to be reassigned. Avoid `var`.
*   **Asynchronous Operations:** Use `async/await` for handling asynchronous operations like API calls and file operations.
*   **Logging:** Use `console.log` for informational messages, `console.warn` for warnings, and `console.error` for errors. Be descriptive in your log messages.
*   **Comments:** Use JSDoc comments for functions and complex logic to explain their purpose, parameters, and return values.

## HTML (UI pages)

*   **Indentation:** Use 2 spaces for indentation.
*   **Structure:** Use semantic HTML5 elements where appropriate.
*   **Attributes:** Quote all attribute values.

## CSS (UI styles)

*   **Indentation:** Use 2 spaces for indentation.
*   **Properties:** Use shorthand properties where applicable (e.g., `margin`, `padding`).
*   **Comments:** Use CSS comments (`/* ... */`) to explain complex styles or sections.
*   **Naming Conventions:** Use BEM (Block, Element, Modifier) or a similar convention for class names to maintain organization.

## Manifest (`manifest.json`)

*   **Formatting:** Use 2 spaces for indentation.
*   **Naming:** Use kebab-case (e.g., `my-extension-id`) for `id` and camelCase (e.g., `myExtensionName`) for `name` and `version`.
*   **Permissions:** Declare all necessary permissions explicitly. Follow the principle of least privilege.

## Naming Convention Summary

*   **File/Directory Names:** kebab-case (e.g., `my-custom-page.html`)
*   **Variable/Function Names:** camelCase (e.g., `myVariable`, `myFunction`)
*   **Class Names (CSS):** BEM (e.g., `block__element--modifier`) or camelCase.
*   **Extension IDs:** kebab-case (e.g., `com.example.myextension`)

## Example Snippet (main.js)

\`\`\`javascript
/**
 * Fetches data from a specified URL using the network API.
 * @param {string} url - The URL to fetch data from.
 * @returns {Promise<object>} - A promise that resolves with the JSON data.
 */
async function fetchData(url) {
  try {
    const response = await Aether.http.get(url);
    if (!response.ok) {
      throw new Error(\`HTTP error! status: \${response.status}\`);
    }
    const data = await response.json();
    console.log(\`Data fetched successfully from \${url}:`, data);
    return data;
  } catch (error) {
    console.error(\`Failed to fetch data from \${url}:`, error);
    throw error; // Re-throw to allow caller to handle
  }
}

// Call the function
fetchData('https://jsonplaceholder.typicode.com/posts/1');
\`\`\`
