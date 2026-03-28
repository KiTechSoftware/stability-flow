FROM squidfunk/mkdocs-material:latest AS build

WORKDIR /site

COPY requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir -r /tmp/requirements.txt

COPY mkdocs.yml /site/mkdocs.yml
COPY docs /site/docs
COPY assets /site/docs/assets

RUN mkdocs build --strict

FROM nginx:alpine

COPY --from=build /site/site /usr/share/nginx/html

EXPOSE 80
