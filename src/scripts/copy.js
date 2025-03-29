// Function to show a notification when text is copied
function showCopyNotification(element) {
  // Create notification element
  const notification = document.createElement('div');
  notification.className = 'copy-notification';
  notification.textContent = 'Copied!';
  
  // Position relative to the clicked button
  const rect = element.getBoundingClientRect();
  notification.style.position = 'fixed';
  notification.style.top = `${rect.top - 30}px`;
  notification.style.left = `${rect.left}px`;
  notification.style.backgroundColor = 'rgba(0, 0, 0, 0.7)';
  notification.style.color = 'white';
  notification.style.padding = '5px 10px';
  notification.style.borderRadius = '4px';
  notification.style.fontSize = '14px';
  notification.style.zIndex = '1000';
  notification.style.opacity = '0';
  notification.style.transition = 'opacity 0.3s ease';
  
  // Add to document
  document.body.appendChild(notification);
  
  // Show with animation
  setTimeout(() => {
    notification.style.opacity = '1';
  }, 10);
  
  // Remove after delay
  setTimeout(() => {
    notification.style.opacity = '0';
    setTimeout(() => {
      document.body.removeChild(notification);
    }, 300);
  }, 1500);
}

// Helper function to extract text content from code blocks
function getCodeText(codeBlock) {
  // For pre>code blocks
  const preCode = codeBlock.querySelector('pre > code');
  if (preCode) {
    return preCode.textContent;
  }
  
  // For chroma/highlight.js style blocks
  const highlightedCode = codeBlock.querySelector('.chroma');
  if (highlightedCode) {
    // Get only text content, ignoring any line numbers
    return Array.from(highlightedCode.querySelectorAll('.line'))
      .map(line => line.textContent)
      .join('\n');
  }
  
  // Fallback: get all text from the element
  return codeBlock.textContent;
}

// Setup event delegation for copy buttons
document.addEventListener('click', function(event) {
  // Handle filename copy button
  const filenameButton = event.target.closest('.copy-filename');
  if (filenameButton) {
    // The onclick handler is already set in HTML
    showCopyNotification(filenameButton);
    return;
  }
  
  // Handle code copy button
  const codeButton = event.target.closest('.code-copy-button');
  if (codeButton) {
    const codeWrapper = codeButton.closest('.code-content-wrapper');
    if (codeWrapper) {
      // Get the code content
      const codeText = getCodeText(codeWrapper);
      
      // Copy to clipboard
      navigator.clipboard.writeText(codeText)
        .then(() => {
          showCopyNotification(codeButton);
        })
        .catch(err => {
          console.error('Could not copy text: ', err);
        });
    }
  }
});