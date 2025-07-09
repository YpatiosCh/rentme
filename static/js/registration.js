let currentStep = 1;
let selectedPlan = null;

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    setupPlanSelection();
    setupNavigation();
    setupTermsCheckbox();
    setupPasswordValidation();
    updateView();
});

// Password Validation
function setupPasswordValidation() {
    const password = document.getElementById('password');
    const confirmPassword = document.getElementById('confirmPassword');
    
    password.addEventListener('input', function() {
        clearFieldError('password');
        if (confirmPassword.value) {
            checkPasswordMatch();
        }
    });
    
    confirmPassword.addEventListener('input', function() {
        checkPasswordMatch();
    });
}

// Check Password Match
function checkPasswordMatch() {
    const password = document.getElementById('password');
    const confirmPassword = document.getElementById('confirmPassword');
    
    clearFieldError('confirmPassword');
    
    if (confirmPassword.value && password.value !== confirmPassword.value) {
        showFieldError('confirmPassword', 'Οι κωδικοί δεν ταιριάζουν');
        return false;
    }
    
    return true;
}

// Plan Selection
function setupPlanSelection() {
    document.querySelectorAll('input[name="plan_selection"]').forEach(radio => {
        radio.addEventListener('change', function() {
            selectedPlan = this.value;
            document.getElementById('selectedPlan').value = selectedPlan;
            updatePlanSummary();
            clearFieldError('plan');
            updateButtons();
        });
    });
}

// Navigation
function setupNavigation() {
    document.getElementById('nextBtn').addEventListener('click', nextStep);
    document.getElementById('prevBtn').addEventListener('click', prevStep);
    document.getElementById('registrationForm').addEventListener('submit', handleSubmit);
}

// Terms Checkbox
function setupTermsCheckbox() {
    document.getElementById('termsAccepted').addEventListener('change', function() {
        clearFieldError('terms');
        updateButtons();
    });
}

// Clear field error
function clearFieldError(fieldName) {
    const errorElement = document.getElementById(fieldName + 'Error');
    if (errorElement) {
        errorElement.textContent = '';
    }
    
    const field = document.getElementById(fieldName) || document.getElementById('termsAccepted');
    if (field) {
        field.classList.remove('error');
    }
}

// Show field error
function showFieldError(fieldName, message) {
    const errorElement = document.getElementById(fieldName + 'Error');
    if (errorElement) {
        errorElement.textContent = message;
    }
    
    const field = document.getElementById(fieldName) || document.getElementById('termsAccepted');
    if (field) {
        field.classList.add('error');
    }
}

// Validate Step 1
function validateStep1() {
    let isValid = true;
    
    // Clear all previous errors
    clearFieldError('email');
    clearFieldError('password');
    clearFieldError('confirmPassword');
    clearFieldError('terms');
    
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const termsAccepted = document.getElementById('termsAccepted').checked;
    
    // Email validation
    if (!email) {
        showFieldError('email', 'Το email είναι υποχρεωτικό');
        isValid = false;
    } else if (!isValidEmail(email)) {
        showFieldError('email', 'Παρακαλώ εισάγετε έγκυρο email');
        isValid = false;
    }
    
    // Password validation
    if (!password) {
        showFieldError('password', 'Ο κωδικός είναι υποχρεωτικός');
        isValid = false;
    } else if (password.length < 8) {
        showFieldError('password', 'Ο κωδικός πρέπει να έχει τουλάχιστον 8 χαρακτήρες');
        isValid = false;
    }
    
    // Confirm password validation
    if (!confirmPassword) {
        showFieldError('confirmPassword', 'Η επιβεβαίωση κωδικού είναι υποχρεωτική');
        isValid = false;
    } else if (password !== confirmPassword) {
        showFieldError('confirmPassword', 'Οι κωδικοί δεν ταιριάζουν');
        isValid = false;
    }
    
    // Terms validation
    if (!termsAccepted) {
        showFieldError('terms', 'Πρέπει να αποδεχτείς τους Όρους Χρήσης και την Πολιτική Απορρήτου');
        isValid = false;
    }
    
    return isValid;
}

// Validate Step 2
function validateStep2() {
    clearFieldError('plan');
    
    if (!selectedPlan) {
        showFieldError('plan', 'Πρέπει να επιλέξεις ένα πλάνο');
        return false;
    }
    
    return true;
}

// Email validation helper
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Next Step
function nextStep() {
    if (currentStep === 1) {
        if (validateStep1()) {
            currentStep = 2;
            updateView();
        }
    }
}

// Previous Step
function prevStep() {
    if (currentStep === 2) {
        currentStep = 1;
        updateView();
    }
}

// Update View
function updateView() {
    // Hide all steps
    document.querySelectorAll('.form-step').forEach(step => {
        step.classList.remove('active');
    });
    
    // Show current step
    document.getElementById(`step${currentStep}`).classList.add('active');
    
    // Update progress bar
    const progressFill = document.getElementById('progressFill');
    progressFill.style.width = (currentStep / 2) * 100 + '%';
    
    // Update step dots
    document.querySelectorAll('.step-dot').forEach((dot, index) => {
        dot.classList.remove('active', 'completed');
        if (index + 1 === currentStep) {
            dot.classList.add('active');
        } else if (index + 1 < currentStep) {
            dot.classList.add('completed');
        }
    });
    
    updateButtons();
}

function updateButtons() {
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');
    const submitBtn = document.getElementById('submitBtn');

    if (currentStep === 1) {
        // Step 1: Email, Password, Όροι
        prevBtn.classList.add('hidden');
        nextBtn.classList.remove('hidden');
        submitBtn.classList.add('hidden');

    } else if (currentStep === 2) {
        // Step 2: Plan Selection
        prevBtn.classList.remove('hidden');
        nextBtn.classList.add('hidden');

        if (selectedPlan) {
            submitBtn.classList.remove('hidden');
            submitBtn.disabled = false;
        } else {
            submitBtn.classList.remove('hidden');
            submitBtn.disabled = true;
        }
    }
}

// Update Plan Summary
function updatePlanSummary() {
    const planNames = {
        'basic': 'Basic Plan',
        'professional': 'Professional Plan',
        'business': 'Business Plan'
    };
    
    const planPrices = {
        'basic': '€10/μήνα',
        'professional': '€20/μήνα',
        'business': '€40/μήνα'
    };
    
    if (selectedPlan) {
        document.getElementById('selectedPlanName').textContent = planNames[selectedPlan];
        document.getElementById('selectedPlanPrice').textContent = planPrices[selectedPlan];
        document.getElementById('totalPrice').textContent = planPrices[selectedPlan];
    }
}

// Handle Form Submit
async function handleSubmit(e) {
    e.preventDefault();
    
    // Final validation
    if (!validateStep2()) {
        return;
    }
    
    const submitBtn = document.getElementById('submitBtn');
    submitBtn.textContent = 'Επεξεργασία...';
    submitBtn.disabled = true;
    
    const formData = {
        plan_id: selectedPlan,
        email: document.getElementById('email').value.trim(),
        password: document.getElementById('password').value,
        terms_accepted: document.getElementById('termsAccepted').checked
    };
    
    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });
        
        const result = await response.json();
        
        if (result.checkout_url) {
            window.location.href = result.checkout_url;
        } else {
            alert('Σφάλμα: ' + (result.error || 'Άγνωστο σφάλμα'));
        }
        
    } catch (error) {
        alert('Σφάλμα δικτύου. Δοκιμάστε ξανά.');
    } finally {
        submitBtn.textContent = 'Πληρωμή τώρα';
        submitBtn.disabled = false;
    }
}

// Modal Functions
function openTermsModal() {
    document.getElementById('termsModal').classList.add('active');
}

function openPrivacyModal() {
    document.getElementById('privacyModal').classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

// Password Toggle Function
function togglePassword(fieldId) {
    const field = document.getElementById(fieldId);
    const button = field.nextElementSibling;
    
    if (field.type === 'password') {
        field.type = 'text';
        button.textContent = '🙈';
    } else {
        field.type = 'password';
        button.textContent = '👁️';
    }
}

// Close modal on outside click
document.addEventListener('click', function(e) {
    if (e.target.classList.contains('modal')) {
        e.target.classList.remove('active');
    }
});