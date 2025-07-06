// Registration Form Controller
class RegistrationForm {
    constructor() {
        this.currentStep = 1;
        this.totalSteps = 5;
        this.selectedPlan = null;
        
        this.init();
    }

    init() {
        this.bindEvents();
        this.updateProgressBar();
        this.setupBusinessNameToggle();
        this.setupPasswordValidation();
    }

    bindEvents() {
        // Plan selection
        document.querySelectorAll('.plan-select-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.selectPlan(e));
        });

        // Navigation buttons
        document.getElementById('nextBtn').addEventListener('click', () => this.nextStep());
        document.getElementById('prevBtn').addEventListener('click', () => this.prevStep());
        
        // Form submission
        document.getElementById('registrationForm').addEventListener('submit', (e) => this.handleSubmit(e));

        // Real-time validation
        this.setupRealTimeValidation();
    }

    setupBusinessNameToggle() {
        const businessNameField = document.getElementById('businessName');
        const businessDescriptionGroup = document.getElementById('businessDescriptionGroup');
        
        businessNameField.addEventListener('input', () => {
            if (businessNameField.value.trim()) {
                businessDescriptionGroup.style.display = 'block';
            } else {
                businessDescriptionGroup.style.display = 'none';
                document.getElementById('businessDescription').value = '';
            }
        });
    }

    setupPasswordValidation() {
        const password = document.getElementById('password');
        const confirmPassword = document.getElementById('confirmPassword');
        
        confirmPassword.addEventListener('input', () => {
            if (password.value !== confirmPassword.value) {
                confirmPassword.setCustomValidity('Οι κωδικοί δεν ταιριάζουν');
            } else {
                confirmPassword.setCustomValidity('');
            }
        });
    }

    setupRealTimeValidation() {
        // Email validation
        const emailField = document.getElementById('email');
        emailField.addEventListener('blur', () => {
            this.validateEmail(emailField);
        });

        // Phone validation
        const phoneField = document.getElementById('phone');
        phoneField.addEventListener('input', () => {
            this.formatPhoneNumber(phoneField);
        });

        // Required field validation
        document.querySelectorAll('input[required], textarea[required]').forEach(field => {
            field.addEventListener('blur', () => {
                this.validateRequiredField(field);
            });
        });
    }

    selectPlan(e) {
        e.preventDefault();
        
        // Remove previous selections
        document.querySelectorAll('.plan-card').forEach(card => {
            card.classList.remove('selected');
        });
        
        // Select current plan
        const planCard = e.target.closest('.plan-card');
        const planId = e.target.dataset.plan;
        
        planCard.classList.add('selected');
        this.selectedPlan = planId;
        
        // Update hidden input
        document.getElementById('selectedPlan').value = planId;
        
        // Update button text
        document.querySelectorAll('.plan-select-btn').forEach(btn => {
            btn.textContent = btn.textContent.replace('✓ Επιλέχθηκε', '').replace('Επιλογή ', 'Επιλογή ');
        });
        
        e.target.textContent = '✓ Επιλέχθηκε';
        
        // Enable next button
        this.validateCurrentStep();
    }

    nextStep() {
        if (this.validateCurrentStep()) {
            if (this.currentStep < this.totalSteps) {
                this.currentStep++;
                this.updateView();
            }
        }
    }

    prevStep() {
        if (this.currentStep > 1) {
            this.currentStep--;
            this.updateView();
        }
    }

    updateView() {
        // Hide all steps
        document.querySelectorAll('.form-step').forEach(step => {
            step.classList.remove('active');
        });
        
        // Show current step
        document.getElementById(`step${this.currentStep}`).classList.add('active');
        
        // Update progress bar
        this.updateProgressBar();
        
        // Update step indicators
        this.updateStepIndicators();
        
        // Update navigation buttons
        this.updateNavigationButtons();
    }

    updateProgressBar() {
        const progressFill = document.getElementById('progressFill');
        const percentage = (this.currentStep / this.totalSteps) * 100;
        progressFill.style.width = `${percentage}%`;
    }

    updateStepIndicators() {
        document.querySelectorAll('.step').forEach((step, index) => {
            const stepNumber = index + 1;
            
            step.classList.remove('active', 'completed');
            
            if (stepNumber === this.currentStep) {
                step.classList.add('active');
            } else if (stepNumber < this.currentStep) {
                step.classList.add('completed');
            }
        });
    }

    updateNavigationButtons() {
        const prevBtn = document.getElementById('prevBtn');
        const nextBtn = document.getElementById('nextBtn');
        const submitBtn = document.getElementById('submitBtn');
        
        // Previous button
        prevBtn.style.display = this.currentStep === 1 ? 'none' : 'block';
        
        // Next/Submit button
        if (this.currentStep === this.totalSteps) {
            nextBtn.style.display = 'none';
            submitBtn.style.display = 'block';
        } else {
            nextBtn.style.display = 'block';
            submitBtn.style.display = 'none';
        }
    }

    validateCurrentStep() {
        switch (this.currentStep) {
            case 1:
                return this.validateStep1();
            case 2:
                return this.validateStep2();
            case 3:
                return this.validateStep3();
            case 4:
                return this.validateStep4();
            case 5:
                return this.validateStep5();
            default:
                return false;
        }
    }

    validateStep1() {
        if (!this.selectedPlan) {
            this.showError('Παρακαλώ επιλέξτε ένα πλάνο συνδρομής');
            return false;
        }
        return true;
    }

    validateStep2() {
        const requiredFields = ['firstName', 'lastName', 'email', 'phone'];
        let isValid = true;
        
        requiredFields.forEach(fieldId => {
            const field = document.getElementById(fieldId);
            if (!this.validateRequiredField(field)) {
                isValid = false;
            }
        });
        
        // Additional email validation
        const emailField = document.getElementById('email');
        if (emailField.value && !this.validateEmail(emailField)) {
            isValid = false;
        }
        
        return isValid;
    }

    validateStep3() {
        // Step 3 is optional, always valid
        return true;
    }

    validateStep4() {
        const requiredFields = ['city', 'region', 'prefecture'];
        let isValid = true;
        
        requiredFields.forEach(fieldId => {
            const field = document.getElementById(fieldId);
            if (!this.validateRequiredField(field)) {
                isValid = false;
            }
        });
        
        return isValid;
    }

    validateStep5() {
        const password = document.getElementById('password');
        const confirmPassword = document.getElementById('confirmPassword');
        let isValid = true;
        
        if (!this.validateRequiredField(password)) {
            isValid = false;
        }
        
        if (!this.validateRequiredField(confirmPassword)) {
            isValid = false;
        }
        
        if (password.value && confirmPassword.value && password.value !== confirmPassword.value) {
            this.showFieldError(confirmPassword, 'Οι κωδικοί δεν ταιριάζουν');
            isValid = false;
        }
        
        if (password.value && password.value.length < 8) {
            this.showFieldError(password, 'Ο κωδικός πρέπει να έχει τουλάχιστον 8 χαρακτήρες');
            isValid = false;
        }
        
        return isValid;
    }

    validateRequiredField(field) {
        if (field.required && !field.value.trim()) {
            this.showFieldError(field, 'Αυτό το πεδίο είναι υποχρεωτικό');
            return false;
        }
        
        this.clearFieldError(field);
        return true;
    }

    validateEmail(field) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (field.value && !emailRegex.test(field.value)) {
            this.showFieldError(field, 'Παρακαλώ εισάγετε έγκυρη διεύθυνση email');
            return false;
        }
        
        this.clearFieldError(field);
        return true;
    }

    formatPhoneNumber(field) {
        // Simple Greek phone formatting
        let value = field.value.replace(/\D/g, '');
        
        if (value.startsWith('30')) {
            value = value.substring(2);
        }
        
        if (value.length <= 10) {
            field.value = value;
        }
    }

    showFieldError(field, message) {
        field.classList.add('error');
        
        // Remove existing error message
        const existingError = field.parentNode.querySelector('.error-message');
        if (existingError) {
            existingError.remove();
        }
        
        // Add new error message
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.textContent = message;
        field.parentNode.appendChild(errorDiv);
    }

    clearFieldError(field) {
        field.classList.remove('error');
        const errorMessage = field.parentNode.querySelector('.error-message');
        if (errorMessage) {
            errorMessage.remove();
        }
    }

    showError(message) {
        // Simple alert for now - can be replaced with better UI
        alert(message);
    }

    handleSubmit(e) {
        e.preventDefault();
        
        if (!this.validateCurrentStep()) {
            return;
        }
        
        // Collect all form data
        const formData = this.collectFormData();
        
        // Submit to backend
        this.submitRegistration(formData);
    }

    collectFormData() {
        const form = document.getElementById('registrationForm');
        const formData = new FormData(form);
        
        // Convert to regular object
        const data = {};
        for (let [key, value] of formData.entries()) {
            data[key] = value;
        }
        
        return data;
    }

    async submitRegistration(data) {
        try {
            // Show loading state
            const submitBtn = document.getElementById('submitBtn');
            const originalText = submitBtn.textContent;
            submitBtn.textContent = 'Επεξεργασία...';
            submitBtn.disabled = true;
            
            const response = await fetch('/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
            
            const result = await response.json();
            
            if (result.checkout_url) {
                // Redirect to Stripe checkout
                window.location.href = result.checkout_url;
            } else {
                throw new Error(result.error || 'Registration failed');
            }
            
        } catch (error) {
            console.error('Registration error:', error);
            this.showError('Σφάλμα κατά την εγγραφή. Παρακαλώ δοκιμάστε ξανά.');
            
            // Reset button
            const submitBtn = document.getElementById('submitBtn');
            submitBtn.textContent = originalText;
            submitBtn.disabled = false;
        }
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new RegistrationForm();
});