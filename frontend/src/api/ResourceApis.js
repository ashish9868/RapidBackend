import axios from "axios";
import { ToastsUtil } from "../utils/ToastsUtil";

axios.defaults.withCredentials = true;


export const BASE_URL = `/api/v1`

const HttpClient = axios.create({
    baseURL: BASE_URL,
    headers: { "Content-type": "application/json" },
});

HttpClient.interceptors.request.use((config) => {
    const contentType = Object.entries(config.headers)?.[1]?.[1];
    config.headers["Content-type"] = contentType || `application/json`;
    // share user's local timezone
    const tz = new Date().getTimezoneOffset()
    if (config.params && typeof config.params === 'object') {
        config.params['tz'] = tz
    } else {
        config.params = { tz }
    }
    return config;
});

HttpClient.interceptors.response.use(
    (response) => {
        const contentDisposition = response?.headers?.["content-disposition"];
        if (contentDisposition && contentDisposition?.includes?.("attachment")) {
            const filename = contentDisposition
                ?.split?.(";")
                ?.find?.((part) => part?.trim()?.startsWith?.("filename"))
                ?.split?.("=")?.[1]
                ?.trim();
            return { data: response.data, filename, success: true };
        }

        return { data: response.data, success: true };
    },
    async (error) => {
        const errorResponse = error?.response;
        switch (errorResponse?.status) {
            case 401:
                ToastsUtil.showError("Session Expired!");
                window.location = '/'
                return Promise.resolve({ success: false, code: 401, errors: { global: 'Unauthorized.' } });

            case 403:
                ToastsUtil.showError("Permission Denied!");
                return Promise.resolve({ success: false, code: 403, errors: { global: 'Permission Denied.' } });
            case 404:
                return Promise.resolve({ success: false, code: 404, errors: { global: 'Record not found.' } });
            case 503:
                ToastsUtil.showError(`Site is undergoing maintaince please try after some time.`)
                return Promise.resolve({ success: false, code: 503, errors: { global: 'Site is undergoing maintaince please try after some time.' } });
            case 422:
                let errors = errorResponse?.data ?? {};
                if (errors instanceof Blob) {
                    try {
                        const text = await errors.text();       // Convert Blob to text
                        errors = JSON.parse(text);              // Parse the text as JSON  
                    } catch (e) {

                    }
                }

                try {
                    ToastsUtil.showError("Please resolve errors to continue.")
                } catch (e) {

                }

                if (errors?.global) {
                    ToastsUtil.showError(errors?.global)
                }
                return Promise.resolve({ success: false, errors, code: 422 });

            case 429:
                ToastsUtil.showError("Too Many Requests!");
                return Promise.resolve({ success: false, code: 429, errors: { global: 'Too many requests.' } });

            case 500:
                ToastsUtil.showError("Something Went Wrong!");
                return Promise.resolve({
                    success: false,
                    code: 500,
                    errors: { global: errorResponse?.data?.message || "Network Error" },
                });

            default:
                ToastsUtil.showError("Something Went Wrong!");
                console.log(error)
                return Promise.reject({ success: false, code: 500, d: 'dddd' });
        }
        // return Promise.reject(error);
    }
);



export const Resources = {
    LOGIN: 'login',
    LOGOUT: 'logout',
    Me: 'me',
    RESET_PASSWORD: 'reset-password',
}


export const buildEndpoint = (resource, config = {
    addBaseUrl: false,
    routeParams: {},
    queryParams: {},
    filters: {},
}) => {
    const routeParams = config?.routeParams || {}
    const queryParams = config?.queryParams || {}
    const filters = config?.filters || {}

    let url = `${config?.addBaseUrl ? BASE_URL : ''}/${resource}${routeParams?.id ? '/:id' : ''}?`
    if (resource) {
        Object.keys(routeParams).forEach(routeParam => {
            url = url.replace(`:${routeParam}`, routeParams[routeParam])
        })
        Object.keys(queryParams).forEach(queryParam => {
            let val = queryParams[queryParam]
            if (Array.isArray(val) || typeof val === 'object') {
                val = JSON.stringify(val)
            }
            url = `${url}&${queryParam}=${encodeURIComponent(val)}`
        })

        Object.keys(filters).forEach(filter => {
            let val = filters[filter]
            if (Array.isArray(val) || typeof val === 'object') {
                val = JSON.stringify(val)
            }
            url = `${url}&${`filters[${filter}]`}=${encodeURIComponent(val)}`
        })
    }
    return url
}

export const ResourceApis = {
    getById: async (resource, id = 0, config = {
        queryParams: {},
        routeParams: {},
        filters: {},
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const routeParams = config?.routeParams || {}
        const endpoint = buildEndpoint(`${resource}`, {
            ...config,
            routeParams: { ...routeParams, id }
        })
        const response = await HttpClient.get(endpoint)
        const data = response?.success ? response.data : null
        callbacks?.onComplete && (await callbacks?.onComplete(response?.success, data, response?.errors))
        return data
    },
    /**
     * Get Paginated Resources
     * 
     * @param {string} resource 
     * @param {{routeParams: Object, queryParams: Object, filters: Object}} config 
     * @returns 
     */
    getPaginated: async (resource, config = {
        routeParams: {},
        queryParams: {},
        filters: {}
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const routeParams = config?.routeParams || {}
        const endpoint = buildEndpoint(`${resource}`, config)
        const response = await HttpClient.get(endpoint)

        if (response?.success) {
            let finalData = null
            if (Array.isArray(response?.data?.results)) {
                finalData = response.data
            } else if (Array.isArray(response?.data)) {
                finalData = {
                    page: 1,
                    from: response.data.length > 0 ? 1 : 0,
                    to: response.data.length,
                    total: response.data.length,
                    results: response.data,
                    total_pages: 1,
                }
            } else {
                const results = [response.data]
                finalData = {
                    page: 1,
                    from: 1,
                    to: results.length,
                    total: results.length,
                    results,
                    total_pages: 1
                }
            }
            callbacks?.onComplete && (await callbacks?.onComplete(true, finalData, response?.errors))
            return finalData
        } else {
            console.error(`Error: `, response?.errors ?? {})
            const finalData = {
                page: 1,
                from: 0,
                to: 1,
                total: 0,
                results: [],
                total_pages: 1
            }
            callbacks?.onComplete && (await callbacks?.onComplete(true, finalData, response?.errors))
            return finalData
        }
    },

    create: async (resource, payload = {}, config = {
        routeParams: {},
        queryParams: {},
        filters: {}
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const routeParams = config?.routeParams || {}
        const url = buildEndpoint(`${resource}`, config)
        const response = await HttpClient.post(url, payload, {
            headers: {
                ...(payload instanceof FormData && {
                    'Content-Type': 'multipart/form-data'
                })
            }
        })
        callbacks?.onComplete && (await callbacks?.onComplete(response?.success, response?.data, response?.errors))
        return response
    },
    update: async (resource, id = 0, payload = {}, config = {
        routeParams: {},
        queryParams: {},
        filters: {}
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const routeParams = config?.routeParams || {}
        const url = buildEndpoint(`${resource}`, {
            ...config,
            routeParams: { ...routeParams, id }
        })

        const response = await HttpClient.put(url, payload, {
            headers: {
                ...(payload instanceof FormData && {
                    'Content-Type': 'multipart/form-data'
                })
            }
        })
        callbacks?.onComplete && (await callbacks?.onComplete(response?.success, response?.data, response?.errors))
        return response
    },
    updatePartial: async (resource, id = 0, payload = {}, config = {
        routeParams: {},
        queryParams: {},
        filters: {}
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const routeParams = config?.routeParams || {}
        const url = buildEndpoint(`${resource}`, {
            ...config,
            routeParams: { ...routeParams, id }
        })

        const response = await HttpClient.patch(url, payload, {
            headers: {
                ...(payload instanceof FormData && {
                    'Content-Type': 'multipart/form-data'
                })
            }
        })
        callbacks?.onComplete && (await callbacks?.onComplete(response?.success, response?.data, response?.errors))
        return response
    },
    delete: async (resource, id = 0, config = {
        queryParams: {},
        routeParams: {},
        filters: {},
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const routeParams = config?.routeParams || {}
        const endpoint = buildEndpoint(`${resource}`, {
            ...config,
            routeParams: { ...routeParams, id }
        })
        const response = await HttpClient.delete(endpoint)
        callbacks?.onComplete && (await callbacks?.onComplete(response?.success, response?.data, response?.errors))
        return response
    },

    download: async (resource, is_post = false, force_download = false, config = {
        queryParams: {},
        routeParams: {},
        filters: {},
    }, callbacks = {
        onBefore: async () => { },
        onComplete: async (success = true, data = null, errors = {}) => { }
    }) => {
        callbacks?.onBefore && (await callbacks?.onBefore())
        const queryParams = config?.queryParams || { download: 1 }
        const endpoint = buildEndpoint(resource, {
            ...config,
            queryParams
        })
        try {
            const response = await HttpClient.request({
                url: endpoint,
                method: is_post ? 'POST' : 'GET',
                responseType: 'blob'
            })
            const blob = new Blob([response.data], { type: response.data.type })
            const fileUrl = window.URL.createObjectURL(blob);
            if (force_download) {
                const a = document.createElement('a')
                a.href = fileUrl
                a.download = response?.filename
                document.body.appendChild(a)
                a.click()
                setTimeout(() => { a.remove() }, 2000)
                window.URL.revokeObjectURL(fileUrl);
                callbacks?.onComplete && (await callbacks?.onComplete(response?.success, response?.data, response?.errors))
            }
            return fileUrl
        } catch (e) {
            console.log("Error bulkDownload", e);
            callbacks?.onComplete && (await callbacks?.onComplete(false, null, {
                global: 'Unable to download the file'
            }))
        }
    }
}