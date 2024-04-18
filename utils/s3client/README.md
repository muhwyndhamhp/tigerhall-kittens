# Using Cloudflare R2 for Image CDN
Uploading image to the server is easy, and it's a common practice to do so. However, serving the image to the user is a different story. Naively using local storage to serve the image is not a good idea for the following reasons:
- The image will be served from the same server as the main server, which can cause the server to be overloaded if the image is requested too many times.
- It can't be hosted on a serverless platform, as the serverless platform usually doesn't have a persistent storage.
- The image can't be cached easily, as the server doesn't have a cache laye. And it will not be cached globally, which can cause the image to be slow to load for users that are far from the server.
- It will incurr Egress Cost for the server to serve the image. Not a good thing to do in serverless platform where egress cost is expensive.

To solve this problem, we can use Cloudflare R2 as the image CDN. Cloudflare R2 is an S3-compatible object storage that is built on top of Cloudflare's global network. By using Cloudflare R2, we can serve the image to the user with the following benefits:
- The image will be cached globally, which can reduce the load on the server and reduce the latency for the user to load the image.
- Now we can deploy our service to a serverless platform, as the image is not stored in the server.
- The image can be served with a custom domain, which can make the image URL more user-friendly.
- Cloudflare R2 is incredibly cheap, and will not incurr Egress Cost a very good fit for serving images. 
