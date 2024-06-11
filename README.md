# Sitemap Builder

A sitemap is basically a map of all of the pages within a specific domain. They are used by search engines and other tools to inform them of all of the pages on your domain.

One way these can be built is by first visiting the root page of the website and making a list of every link on that page that goes to a page on the same domain. For instance, on `google.com` you might find a link to `google.com/search/` along with several other links.

Once you have created the list of links, you could then visit each and add any new links to your list. By repeating this step over and over you would eventually visit every page that on the domain that can be reached by following links from the root page.

sitemap builder outputs the data in the following XML format:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
  </url>
  <url>
    <loc>http://www.example.com/dogs</loc>
  </url>
</urlset>
```

_Note: This should be the same as the [standard sitemap protocol](https://www.sitemaps.org/index.html)_

Where each page is listed in its own `<url>` tag and includes the `<loc>` tag inside of it.

In order to complete this exercise I highly recommend first doing the [link parser exercise](https://github.com/gophercises/link) and using the package created in it to parse your HTML pages for links.

From there you will likely need to figure out a way to determine if a link goes to the same domain or a different one. If it goes to a different domain we shouldn't include it in our sitemap builder, but if it is the same domain we should. Remember that links to the same domain can be in the format of `/just-the-path` or `https://domain.com/with-domain`, but both go to the same domain.
