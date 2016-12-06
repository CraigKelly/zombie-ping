// Taken from https://github.com/mozilla/serviceworker-cookbook/

// Register event listener for the 'push' event.
self.addEventListener('push', function(event) {
    // Keep the service worker alive until the notification is created.
    event.waitUntil(
        self.registration.showNotification('Zombie-Ping found a problem!', {
            tag: 'unique-id-for-this-notification',
            body: 'There is an issue with at least one URL being watched by Zombie-Ping with a problem',
            vibration: [50, 50, 100, 50, 250, 50, 100, 50, 50],
            actions: [
                { action: 'open', title: 'Open the thing' }
            ]
        })
    );
});

// Event for user clicking on action
self.addEventListener('notificationclick', function(evt) {
    console.log("Notification clicked!");
    if (evt.action === 'open') {
        console.log("Open the correct page!");
        alert('Need to open the correct page');
    }
    else {
        console.log('Unknown notification clicked:', evt, evt.action || 'noaction');
    }
    evt.notification.close();
}, false);
