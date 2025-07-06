function selectPlan(planId) {
    document.getElementById('selectedPlan').value = planId;
    document.getElementById('subscriptionForm').style.display = 'block';
    document.getElementById('subscriptionForm').scrollIntoView();
}

document.getElementById('checkoutForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const data = {
        plan_id: formData.get('plan_id'),
        first_name: formData.get('first_name'),
        last_name: formData.get('last_name'),
        email: formData.get('email'),
        phone: formData.get('phone'),
        business_name: formData.get('business_name'),
        address: formData.get('address'),
        city: formData.get('city'),
        region: formData.get('region'),
        postal_code: formData.get('postal_code')
    };

    try {
        const response = await fetch('/checkout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        const result = await response.json();
        
        if (result.checkout_url) {
            window.location.href = result.checkout_url;
        } else {
            alert('Error creating checkout session');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error creating checkout session');
    }
});