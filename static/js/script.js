// Port selection handling
document.addEventListener('DOMContentLoaded', function() {
    // Hide loading overlay on page load
    document.getElementById('loadingOverlay').style.display = 'none';

    // Show/hide custom range inputs
    document.getElementById('customRange').addEventListener('change', function() {
        document.getElementById('customRangeInputs').style.display = this.checked ? 'flex' : 'none';
        if (this.checked) {
            document.getElementById('specificPortsInput').style.display = 'none';
        }
    });

    // Show/hide specific ports input
    document.getElementById('specificPorts').addEventListener('change', function() {
        document.getElementById('specificPortsInput').style.display = this.checked ? 'block' : 'none';
        if (this.checked) {
            document.getElementById('customRangeInputs').style.display = 'none';
        }
    });

    // Hide inputs when other options are selected
    document.getElementById('commonPorts').addEventListener('change', function() {
        if (this.checked) {
            document.getElementById('customRangeInputs').style.display = 'none';
            document.getElementById('specificPortsInput').style.display = 'none';
        }
    });

    document.getElementById('wellKnownPorts').addEventListener('change', function() {
        if (this.checked) {
            document.getElementById('customRangeInputs').style.display = 'none';
            document.getElementById('specificPortsInput').style.display = 'none';
        }
    });

    // Form submission handling
    document.getElementById('scanForm').onsubmit = function(e) {
        // Get the IP address value
        const ipAddress = document.getElementById('ip').value.trim();
        
        // Basic IP address validation
        const ipPattern = /^(\d{1,3}\.){3}\d{1,3}$/;
        if (!ipPattern.test(ipAddress)) {
            e.preventDefault(); // Prevent form submission
            showNotification('Please enter a valid IP address', 'error');
            return false;
        }

        // Validate custom range if selected
        if (document.getElementById('customRange').checked) {
            const startPort = parseInt(document.getElementById('startPort').value);
            const endPort = parseInt(document.getElementById('endPort').value);
            
            if (isNaN(startPort) || isNaN(endPort) || startPort < 1 || endPort > 65535 || startPort > endPort) {
                e.preventDefault(); // Prevent form submission
                showNotification('Please enter valid port range (1-65535)', 'error');
                return false;
            }
        }

        // Validate specific ports if selected
        if (document.getElementById('specificPorts').checked) {
            const portList = document.getElementById('portList').value.trim();
            if (!portList) {
                e.preventDefault(); // Prevent form submission
                showNotification('Please enter specific ports', 'error');
                return false;
            }

            // Validate each port number
            const ports = portList.split(',').map(p => parseInt(p.trim()));
            for (let port of ports) {
                if (isNaN(port) || port < 1 || port > 65535) {
                    e.preventDefault(); // Prevent form submission
                    showNotification('Invalid port number detected. Ports must be between 1-65535', 'error');
                    return false;
                }
            }
        }

        // If all validations pass, show loading overlay
        document.getElementById('loadingOverlay').style.display = 'flex';
        document.getElementById('scanButton').disabled = true;
        document.getElementById('scanButton').innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>Scanning...';
        return true;
    };
});

// Results copying functionality
function copyResults() {
    const resultsText = document.getElementById('resultsBox').innerText;
    navigator.clipboard.writeText(resultsText)
        .then(function() {
            showNotification('Results copied to clipboard!', 'success');
        })
        .catch(function(err) {
            console.error('Failed to copy results:', err);
            showNotification('Failed to copy results', 'error');
        });
}

// Notification system
function showNotification(message, type) {
    const notification = document.createElement('div');
    notification.className = `alert alert-${type === 'success' ? 'success' : 'danger'} notification`;
    notification.style.position = 'fixed';
    notification.style.top = '20px';
    notification.style.right = '20px';
    notification.style.zIndex = '1050';
    notification.style.padding = '10px 20px';
    notification.style.borderRadius = '8px';
    notification.style.animation = 'fadeInOut 3s forwards';
    notification.style.maxWidth = '90%';
    notification.style.width = '400px';
    notification.innerHTML = `
        <i class="fas fa-${type === 'success' ? 'check-circle' : 'exclamation-circle'} me-2"></i>
        ${message}
    `;

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.remove();
    }, 3000);
}

// Add custom CSS animation
const style = document.createElement('style');
style.textContent = `
    @keyframes fadeInOut {
        0% { opacity: 0; transform: translateY(-20px); }
        10% { opacity: 1; transform: translateY(0); }
        90% { opacity: 1; transform: translateY(0); }
        100% { opacity: 0; transform: translateY(-20px); }
    }
`;
document.head.appendChild(style); 