<script>
import CommonMixins from '../mixins/CommonMixins'
import { Toast } from 'bootstrap'
import { mailbox } from '../stores/mailbox'
import { pagination } from '../stores/pagination'

export default {
	mixins: [CommonMixins],

	data() {
		return {
			pagination,
			mailbox,
			toastMessage: false,
			reconnectRefresh: false,
			socketURI: false,
			pauseNotifications: false, // prevent spamming
			version: false
		}
	},

	mounted() {
		let d = document.getElementById('app')
		if (d) {
			this.version = d.dataset.version
		}

		let windowUrl = window.location;
		this.sseURI = new URL(windowUrl .protocol + "//" + windowUrl.host + "/.well-known/mercure" );
		this.sseURI.searchParams.append('topic', 'mail');
		this.connect()
		
		mailbox.notificationsSupported = window.isSecureContext
			&& ("Notification" in window && Notification.permission !== "denied")
		mailbox.notificationsEnabled = mailbox.notificationsSupported && Notification.permission == "granted"
	},

	methods: {
		// see connect
		connect: function () {
			const sse = new EventSource(this.sseURI, { withCredentials: true });
			let self = this

			// The callback will be called every time an update is published on the Mercure SSE stream
			sse.onmessage = function (e) {
				let response
				try {
					response = JSON.parse(e.data)
				} catch (e) {
					return
				}

				console.log(response);

				// new messages
				if (response.state == "received") {
					console.log("received new message")
					if (!mailbox.searching) {
						if (pagination.start < 1) {
							// push results directly into first page
							mailbox.messages.unshift(response.Data)
							if (mailbox.messages.length > pagination.limit) {
								mailbox.messages.pop()
							}
						} else {
							// update pagination offset
							pagination.start++
						}
					}

					for (let i in response.Data.Tags) {
						if (mailbox.tags.indexOf(response.Data.Tags[i]) < 0) {
							mailbox.tags.push(response.Data.Tags[i])
							mailbox.tags.sort()
						}
					}

					// send notifications
					if (!self.pauseNotifications) {
						self.pauseNotifications = true
						let from = response.Data.From != null ? response.Data.From : '[unknown]'
						self.browserNotify("New mail from: " + from, response.Data.Subject)
						self.setMessageToast(response.Data)
						// delay notifications by 2s
						window.setTimeout(() => { self.pauseNotifications = false }, 2000)
					}
				} else if (response.state  == "processed") {
					// send notifications
					if (!self.pauseNotifications) {
						self.pauseNotifications = true
						let from = response.Data.From != null ? response.Data.From : '[unknown]'
						self.browserNotify("Processed mail: " + from, response.Data.Subject)
						self.setMessageToast(response.Data)
						// delay notifications by 2s
						window.setTimeout(() => { self.pauseNotifications = false }, 2000)
					}
				} else if (response.state == "sent" ) {
					// send notifications
					if (!self.pauseNotifications) {
						self.pauseNotifications = true
						let to = response.Data.To != null ? response.Data.To : '[unknown]'
						self.browserNotify("Sent mail to: " + to, response.Data.Subject)
						self.setMessageToast(response.Data)
						// delay notifications by 2s
						window.setTimeout(() => { self.pauseNotifications = false }, 2000)
					}
				} else if (response.state == "tracked") {
					// send notifications
					if (!self.pauseNotifications) {
						self.pauseNotifications = true
						let to = response.Data.To != null ? response.Data.To : '[unknown]'
						self.browserNotify("Mail read: " + to, response.Data.Subject)
						self.setMessageToast(response.Data)
						// delay notifications by 2s
						window.setTimeout(() => { self.pauseNotifications = false }, 2000)
					}
				}
			};
		},

		browserNotify: function (title, message) {
			if (!("Notification" in window)) {
				return
			}

			if (Notification.permission === "granted") {
				let b = message.Subject
				let options = {
					body: message,
					icon: this.resolve('/notification.png')
				}
				new Notification(title, options)
			}
		},

		setMessageToast: function (m) {
			// don't display if browser notifications are enabled, or a toast is already displayed
			if (mailbox.notificationsEnabled || this.toastMessage) {
				return
			}

			this.toastMessage = m

			let self = this
			let el = document.getElementById('messageToast')
			if (el) {
				el.addEventListener('hidden.bs.toast', () => {
					self.toastMessage = false
				})

				Toast.getOrCreateInstance(el).show()
			}
		},

		closeToast: function () {
			let el = document.getElementById('messageToast')
			if (el) {
				Toast.getOrCreateInstance(el).hide()
			}
		},
	},
}
</script>

<template>
	<div class="toast-container position-fixed bottom-0 end-0 p-3">
		<div id="messageToast" class="toast" role="alert" aria-live="assertive" aria-atomic="true">
			<div class="toast-header" v-if="toastMessage">
				<i class="bi bi-envelope-exclamation-fill me-2"></i>
				<strong class="me-auto">
					<RouterLink :to="'/view/' + toastMessage.ID" @click="closeToast">New message</RouterLink>
				</strong>
				<button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
			</div>

			<div class="toast-body">
				<div>
					<RouterLink :to="'/view/' + toastMessage.ID" class="d-block text-truncate text-body-secondary"
						@click="closeToast">
						<template v-if="toastMessage.Subject != ''">{{ toastMessage.Subject }}</template>
						<template v-else>
							[ no subject ]
						</template>
					</RouterLink>
				</div>
			</div>
		</div>
	</div>
</template>
