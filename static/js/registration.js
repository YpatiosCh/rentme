let currentStep = 1;
let selectedPlan = null;
let stripe;
let elements;
let paymentElement;

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    setupPlanSelection();
    setupNavigation();
    setupTermsCheckbox();
    setupPasswordValidation();
    updateView();
});

// Initialize Stripe
async function initializeStripe() {
    try {
        // Get publishable key from backend
        const response = await fetch('/stripe/config');
        const config = await response.json();
        
        if (config.publishable_key) {
            stripe = Stripe(config.publishable_key);
            elements = stripe.elements();
        } else {
            console.error('Failed to get Stripe publishable key');
        }
    } catch (error) {
        console.error('Error initializing Stripe:', error);
    }
    
    // Create card element
    cardElement = elements.create('card', {
        style: {
            base: {
                fontSize: '15px',
                color: '#2d3748',
                fontFamily: '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
                '::placeholder': {
                    color: '#a0aec0',
                },
            },
            invalid: {
                color: '#e53e3e',
            },
        },
    });
}

// Mount card elements when Step 3 becomes active
function mountCardElements() {
    if (cardNumberElement && !cardNumberElement._mounted) {
        try {
            cardNumberElement.mount('#card-number-element');
            cardExpiryElement.mount('#card-expiry-element');
            cardCvcElement.mount('#card-cvc-element');
            
            cardNumberElement._mounted = true;
            
            console.log('Card elements mounted successfully');
            
            // Handle real-time validation errors
            cardNumberElement.addEventListener('change', function(event) {
                handleCardError(event, 'card-number-errors', 'card-number-element');
            });
            
            cardExpiryElement.addEventListener('change', function(event) {
                handleCardError(event, 'card-expiry-errors', 'card-expiry-element');
            });
            
            cardCvcElement.addEventListener('change', function(event) {
                handleCardError(event, 'card-cvc-errors', 'card-cvc-element');
            });
        } catch (error) {
            console.error('Error mounting card elements:', error);
        }
    } else {
        console.log('Card elements already mounted or not initialized');
    }
}

function handleCardError(event, errorElementId, inputElementId) {
    const displayError = document.getElementById(errorElementId);
    const cardInput = document.getElementById(inputElementId);
    
    if (event.error) {
        displayError.textContent = event.error.message;
        cardInput.classList.add('error');
    } else {
        displayError.textContent = '';
        cardInput.classList.remove('error');
    }
}

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

function setupNavigation() {
    document.getElementById('nextBtn').addEventListener('click', nextStep);
    document.getElementById('prevBtn').addEventListener('click', prevStep);
    document.getElementById('submitBtn').addEventListener('click', handleSubmit);
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
    
    clearFieldError('email');
    clearFieldError('password');
    clearFieldError('confirmPassword');
    clearFieldError('terms');
    
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const termsAccepted = document.getElementById('termsAccepted').checked;
    
    if (!email) {
        showFieldError('email', 'Το email είναι υποχρεωτικό');
        isValid = false;
    } else if (!isValidEmail(email)) {
        showFieldError('email', 'Παρακαλώ εισάγετε έγκυρο email');
        isValid = false;
    }
    
    if (!password) {
        showFieldError('password', 'Ο κωδικός είναι υποχρεωτικός');
        isValid = false;
    } else if (password.length < 8) {
        showFieldError('password', 'Ο κωδικός πρέπει να έχει τουλάχιστον 8 χαρακτήρες');
        isValid = false;
    }
    
    if (!confirmPassword) {
        showFieldError('confirmPassword', 'Η επιβεβαίωση κωδικού είναι υποχρεωτική');
        isValid = false;
    } else if (password !== confirmPassword) {
        showFieldError('confirmPassword', 'Οι κωδικοί δεν ταιριάζουν');
        isValid = false;
    }
    
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
    } else if (currentStep === 2) {
        if (validateStep2()) {
            currentStep = 3;
            updateView();
            updateFinalPlanSummary();
            initializePayment(); // Initialize payment when entering Step 3
        }
    }
}

// Previous Step
function prevStep() {
    if (currentStep === 3) {
        currentStep = 2;
        updateView();
    } else if (currentStep === 2) {
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
    progressFill.style.width = (currentStep / 3) * 100 + '%';
    
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
        prevBtn.classList.add('hidden');
        nextBtn.classList.remove('hidden');
        submitBtn.classList.add('hidden');
    } else if (currentStep === 2) {
        prevBtn.classList.remove('hidden');
        nextBtn.classList.remove('hidden');
        submitBtn.classList.add('hidden');
    } else if (currentStep === 3) {
        prevBtn.classList.remove('hidden');
        nextBtn.classList.add('hidden');
        submitBtn.classList.remove('hidden');
        submitBtn.disabled = false;
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

// Update Final Plan Summary for Step 3
function updateFinalPlanSummary() {
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
        document.getElementById('finalPlanName').textContent = planNames[selectedPlan];
        document.getElementById('finalPlanPrice').textContent = planPrices[selectedPlan];
    }
}
// Handle Form Submit
async function handleSubmit(e) {
    e.preventDefault();
    
    const submitBtn = document.getElementById('submitBtn');
    submitBtn.textContent = 'Επεξεργασία...';
    submitBtn.disabled = true;
    
    try {
        // Confirm payment for subscription
        const { error, paymentIntent } = await stripe.confirmPayment({
            elements,
            redirect: 'if_required',
        });

        if (error) {
            showPopup('error', 'Σφάλμα πληρωμής', error.message);
        } else if (paymentIntent && paymentIntent.status === 'succeeded') {
            // Payment successful - complete registration with subscription data
            try {
                const registrationData = {
                    email: document.getElementById('email').value.trim(),
                    password: document.getElementById('password').value,
                    plan_id: selectedPlan,
                    payment_intent_id: paymentIntent.id,
                    subscription_id: window.subscriptionData.subscriptionId,
                    customer_id: window.subscriptionData.customerId
                };
                
                const regResponse = await fetch('/complete-registration', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(registrationData)
                });
                
                if (regResponse.ok) {
                    showPopup('success', 'Επιτυχία!', 'Η εγγραφή σας ολοκληρώθηκε επιτυχώς! Το subscription σας είναι ενεργό.', true);
                } else {
                    const errorData = await regResponse.json();
                    showPopup('error', 'Σφάλμα εγγραφής', 
                        errorData.message || 'Η πληρωμή ήταν επιτυχής αλλά υπήρξε πρόβλημα με την εγγραφή');
                }
                
            } catch (registrationError) {
                console.error('Registration error:', registrationError);
                showPopup('error', 'Σφάλμα εγγραφής', 
                    'Η πληρωμή ήταν επιτυχής αλλά υπήρξε πρόβλημα με την εγγραφή. Επικοινωνήστε μαζί μας.');
            }
        } else {
            showPopup('error', 'Σφάλμα πληρωμής', 'Η πληρωμή δεν ολοκληρώθηκε. Δοκιμάστε ξανά.');
        }
        
    } catch (error) {
        console.error('Payment error:', error);
        showPopup('error', 'Σφάλμα δικτύου', 'Παρουσιάστηκε σφάλμα. Δοκιμάστε ξανά.');
    } finally {
        submitBtn.textContent = 'Πληρωμή τώρα';
        submitBtn.disabled = false;
    }
}

// Show popup with message
function showPopup(type, title, message, withCountdown = false) {
    const popup = document.getElementById('messagePopup');
    const icon = document.getElementById('popupIcon');
    const titleEl = document.getElementById('popupTitle');
    const messageEl = document.getElementById('popupMessage');
    const countdown = document.getElementById('popupCountdown');
    const closeBtn = document.getElementById('popupCloseBtn');
    
    // Set content
    titleEl.textContent = title;
    messageEl.textContent = message;
    
    // Set type
    popup.className = `popup-overlay ${type}`;
    
    if (type === 'success') {
        icon.textContent = '✅';
    } else {
        icon.textContent = '❌';
    }
    
    // Show popup
    popup.classList.remove('hidden');
    
    // Handle countdown and redirect
    if (withCountdown) {
        countdown.classList.remove('hidden');
        closeBtn.style.display = 'none';
        
        let timeLeft = 4;
        const countdownNumber = document.getElementById('countdownNumber');
        
        const timer = setInterval(() => {
            timeLeft--;
            countdownNumber.textContent = timeLeft;
            
            if (timeLeft <= 0) {
                clearInterval(timer);
                window.location.href = '/';
            }
        }, 1000);
    } else {
        countdown.classList.add('hidden');
        closeBtn.style.display = 'inline-flex';
        
        closeBtn.onclick = function() {
            popup.classList.add('hidden');
        };
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
    
    if (e.target.classList.contains('popup-overlay')) {
        e.target.classList.add('hidden');
    }
});
// Initialize Stripe and get clientSecret when Step 3 loads
async function initializePayment() {
    try {
        // Get publishable key
        const configResponse = await fetch('/stripe/config');
        const config = await configResponse.json();
        
        if (!config.publishable_key) {
            console.error('Failed to get Stripe publishable key');
            return;
        }
        
        stripe = Stripe(config.publishable_key);
        
        // Create Subscription
        const paymentResponse = await fetch('/create-subscription', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                plan_id: selectedPlan,
                email: document.getElementById('email').value.trim()
            })
        });
        
        const { clientSecret, subscriptionId, customerId } = await paymentResponse.json();
        console.log('Got subscription data:', { clientSecret, subscriptionId, customerId });
        
        // Store subscription data for later use
        window.subscriptionData = { subscriptionId, customerId };
        
        // Initialize Elements
        elements = stripe.elements({ clientSecret });
        
        // Create Payment Element
        paymentElement = elements.create('payment');
        paymentElement.mount('#payment-element');
        
        console.log('Subscription payment initialized successfully');
        
    } catch (error) {
        console.error('Error initializing payment:', error);
    }
}