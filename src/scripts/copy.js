function showCopyNotification(element) {
  const notification = document.createElement('div');
  notification.className = 'copy-notification';
  notification.textContent = 'Copied!';
  
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
  
  document.body.appendChild(notification);
  
  setTimeout(() => {
    notification.style.opacity = '1';
  }, 10);
  
  setTimeout(() => {
    notification.style.opacity = '0';
    setTimeout(() => {
      document.body.removeChild(notification);
    }, 300);
  }, 1500);
}

function getCodeText(codeBlock) {
  const preCode = codeBlock.querySelector('pre > code');
  if (preCode) {
    return preCode.textContent;
  }
  
  const highlightedCode = codeBlock.querySelector('.chroma');
  if (highlightedCode) {
    return Array.from(highlightedCode.querySelectorAll('.line'))
      .map(line => line.textContent)
      .join('\n');
  }
  
  return codeBlock.textContent;
}

document.addEventListener('click', function(event) {
  const filenameButton = event.target.closest('.copy-filename');
  if (filenameButton) {
    showCopyNotification(filenameButton);
    return;
  }
  
  const codeButton = event.target.closest('.code-copy-button');
  if (codeButton) {
    const codeWrapper = codeButton.closest('.code-content-wrapper');
    if (codeWrapper) {
      const codeText = getCodeText(codeWrapper);
      
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