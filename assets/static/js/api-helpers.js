/**
 * API Helper functions for SSG
 * Centralizes API calls and automatically adds required headers
 *
 * Usage examples:
 *
 * GET request:
 *   apiFetch(`/api/v1/ssg/contents/${contentId}/images`)
 *
 * POST request:
 *   apiFetch(`/api/v1/ssg/contents/${contentId}/images`, {
 *     method: 'POST',
 *     body: formData
 *   })
 *
 * DELETE request:
 *   apiFetch(`/api/v1/ssg/contents/${contentId}/images/delete`, {
 *     method: 'DELETE',
 *     headers: { 'Content-Type': 'application/json' },
 *     body: JSON.stringify({ image_path: path })
 *   })
 *
 * PUT request:
 *   apiFetch(`/api/v1/ssg/contents/${contentId}`, {
 *     method: 'PUT',
 *     headers: { 'Content-Type': 'application/json' },
 *     body: JSON.stringify(data)
 *   })
 */

/**
 * Wrapper around fetch that automatically adds X-Site-Slug header
 * Works with all HTTP methods: GET, POST, PUT, DELETE, PATCH, etc.
 * @param {string} url - The URL to fetch
 * @param {RequestInit} options - Fetch options (method, body, headers, etc.)
 * @returns {Promise<Response>} The fetch response
 */
function apiFetch(url, options = {}) {
  const siteSlug = getSiteSlug();

  const defaultOptions = {
    headers: {
      'X-Site-Slug': siteSlug,
      ...(options.headers || {})
    }
  };

  const mergedOptions = {
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...options.headers
    }
  };

  return fetch(url, mergedOptions);
}

/**
 * Get the site slug from the page context
 * Priority: data attribute > cookie > default
 * @returns {string} The site slug
 */
function getSiteSlug() {
  // Try to get from data attribute first
  const siteSlug = document.body.dataset.siteSlug;
  if (siteSlug) {
    return siteSlug;
  }

  // Try to get from cookie
  const cookies = document.cookie.split(';');
  for (let cookie of cookies) {
    const [name, value] = cookie.trim().split('=');
    if (name === 'site_slug') {
      return value;
    }
  }

  // Default fallback
  return 'structured';
}

/**
 * Get the API base URL
 * TODO: Make this configurable instead of hardcoded
 * @returns {string} The API base URL
 */
function getAPIBaseURL() {
  return 'http://localhost:8081/api/v1';
}
