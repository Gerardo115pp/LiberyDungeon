import { base_domain, downloads_server } from "../services";
import { DownloadProgress } from "@models/Downloads";

export class DownloadProgressTransmisor {
    constructor(download_uuid) {
        this.download_uuid = download_uuid;
        this.host = `wss://${base_domain}${downloads_server}/ws/download-progress?download_uuid=${this.download_uuid}`;
        this.socket = null;

        this.download_progress_callback = null;
        this.download_completed_callback = null;
    }

    connect = () => {
        this.socket = new WebSocket(this.host);
        this.socket.onopen = this.onOpen;
        this.socket.onmessage = this.onMessage;
        this.socket.onclose = this.onClose;
        this.socket.onerror = this.onError;
    }

    disconnect = () => {
        this.socket.close();
    }

    onMessage = message => {
        console.debug("received message from server: ", message)
        const data = JSON.parse(message.data);

        let download_progress;

        try {
            download_progress = new DownloadProgress(data);
        } catch (error) {
            console.error(error);
        }

        if (this.download_progress_callback !== null) {
            this.download_progress_callback(download_progress);
        }   
    }

    onOpen = () => {}

    onClose = () => {}

    onError = error => {}
}