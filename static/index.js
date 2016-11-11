(function(module) {
    // Our endpoint
    var endpoint = null;

    // Simple subscription
    function do_subscribe() {
        navigator.serviceWorker.register('service-worker.js')
            .then(function(registration) {
                // Use the PushManager to get the user's subscription to the push service.
                return registration.pushManager.getSubscription()
                    .then(function(subscription) {
                        // If a subscription was found, return it.
                        if (subscription) {
                            console.log("Starting a subscription!");
                            return subscription;
                        }

                        // Otherwise, subscribe the user (userVisibleOnly allows to specify that we don't plan to
                        // send notifications that don't have a visible effect for the user).
                        console.log("Subscribing");
                        return registration.pushManager.subscribe({
                            userVisibleOnly: true
                        });
                    });
            }).then(function(subscription) {
                console.log("Have a subscription: telling the server about it");
                endpoint = subscription.endpoint;

                // Show curl command to send the notification on the page.
                $('#curl').text('curl -H "TTL: 60" -X POST ' + endpoint);

                // Send the subscription details to the server using the Fetch API.
                fetch('/subscription-register', {
                    method: 'post',
                    headers: {
                        'Content-type': 'application/json'
                    },
                    body: JSON.stringify({
                        endpoint: subscription.endpoint,
                    }),
                });
            });
    }

    // Register a Service Worker if we have permission
    module.init_notifications = function() {
        console.log("Starting notification subscription logic");
        if (Notification.permission !== "granted") {
            console.log("Getting notification permission");
            Notification.requestPermission(function(permission) {
                if (permission === "granted") {
                    console.log("Notification permission granted: subscribing");
                    do_subscribe();
                }
                else {
                    console.log("DENIED!");
                }
            });
        }
        else {
            console.log("We have notification permission - proceeding");
            do_subscribe();
        }
    };
})(window || this);
