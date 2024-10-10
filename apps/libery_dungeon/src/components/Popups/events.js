const events = {
    NEW_NOTIFICATION: 'notification_triggered'
};

export const emitEvent = (eventName, data) => {
    const new_event = new CustomEvent(eventName, {...data});
    document.dispatchEvent(new_event);
}


export const newNotification = message => {
    emitEvent(events.NEW_NOTIFICATION, {detail: {message}});
}

export default events;