const config = useRuntimeConfig()
const baseApi = config.API_BASE_URL
export default defineEventHandler(async (event) => {
	const paths = event.context.params?._
	const method = event.method

	// Handle swagger static assets with proper content-type headers
	if (paths && typeof paths === 'string' && paths.includes('swagger/')) {
		const url = `${baseApi}/api/${paths}`

		try {
			const response = await $fetch.raw(url, {
				method,
				query: getQuery(event),
			})

			// Set appropriate content-type for swagger assets
			if (paths.endsWith('.css')) {
				setHeader(event, 'content-type', 'text/css')
			} else if (paths.endsWith('.js')) {
				setHeader(event, 'content-type', 'application/javascript')
			} else if (paths.endsWith('.png')) {
				setHeader(event, 'content-type', 'image/png')
			} else if (paths.endsWith('.svg')) {
				setHeader(event, 'content-type', 'image/svg+xml')
			}

			// Copy other relevant headers from the response
			const contentType = response.headers.get('content-type')
			if (contentType) {
				setHeader(event, 'content-type', contentType)
			}

			return response._data
		} catch (error) {
			// Fallback to regular proxy if direct fetch fails
		}
	}

	// Only read body for methods that typically have a body
	const hasBody = ['POST', 'PUT', 'PATCH', 'DELETE'].includes(method)

	const options: any = {
		method,
		query: getQuery(event),
	}

	if (hasBody) {
		try {
			// Check if this is a multipart/form-data request (file upload)
			const contentType = getHeader(event, 'content-type')

			if (contentType && contentType.includes('multipart/form-data')) {
				// For file uploads, use proxyRequest to forward the raw stream directly
				return proxyRequest(event, `${baseApi}/api/${paths}`, {
					...options,
					streamRequest: true, // Stream the request body directly
				})
			} else {
				// For regular JSON requests
				options.body = await readBody(event)
			}
		} catch (error) {
			// Handle cases where body is expected but empty or malformed
			options.body = undefined
		}
	}

	return await event.$fetch(`${baseApi}/api/${paths}`, options);
});
