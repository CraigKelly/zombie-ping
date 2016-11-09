// Taken from https://github.com/mozilla/serviceworker-cookbook/

// Register event listener for the 'push' event.
self.addEventListener('push', function(event) {
    // Keep the service worker alive until the notification is created.
    event.waitUntil(
        self.registration.showNotification('Hello There!', {
            tag: 'unique-id-for-this-notification',
            body: 'Alea iacta est',
            vibration: [50, 50, 100, 50, 250, 50, 100, 50, 50],
            actions: [
                { action: 'open', title: 'Open the thing' }
            ]
        })
    );
});

// Event for user clicking on action
self.addEventListener('notificationclick', function(evt) {
    evt.notification.close();
    if (evt.action === 'open') {
        alert('Need to open the correct page');
    }
    else {
        console.log('Unknown notification clicked:', evt);
    }
}, false);
