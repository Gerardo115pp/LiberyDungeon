package models

import "fmt"

type downloadQueueItem struct {
	Download *DownloadRequest
	next     *downloadQueueItem
}

type DownloadQueue struct {
	head *downloadQueueItem
	tail *downloadQueueItem
	len  int
}

func (dq *DownloadQueue) Clear() {
	dq.head = nil
	dq.tail = nil
	dq.len = 0
}

func (dq *DownloadQueue) Enqueue(download *DownloadRequest) {
	var new_item *downloadQueueItem = new(downloadQueueItem)
	new_item.Download = download
	new_item.next = dq.tail
	dq.tail = new_item

	if dq.head == nil {
		dq.head = dq.tail
	}

	dq.len++
	return
}

func (dq *DownloadQueue) Dequeue() (*DownloadRequest, error) {
	if dq.len == 0 {
		return nil, nil
	}

	var download_node *downloadQueueItem = dq.tail

	if !download_node.Download.IsDownloaded() {
		return nil, fmt.Errorf("Download is not completed")
	}

	dq.tail = download_node.next

	if dq.head == download_node {
		dq.head = nil
	}

	dq.len--

	return download_node.Download, nil
}

func (dq *DownloadQueue) Peek() *DownloadRequest {
	if dq.len == 0 {
		return nil
	}

	return dq.tail.Download
}

func (dq *DownloadQueue) Len() int {
	return dq.len
}

func (dq *DownloadQueue) IsEmpty() bool {
	return dq.tail == nil && dq.head == nil
}
