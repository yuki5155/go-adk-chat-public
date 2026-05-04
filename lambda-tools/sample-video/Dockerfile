FROM public.ecr.aws/lambda/python:3.12

RUN dnf install -y tar xz && \
    curl -L https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz \
      -o /tmp/ffmpeg.tar.xz && \
    tar -xf /tmp/ffmpeg.tar.xz -C /tmp && \
    mv /tmp/ffmpeg-*-amd64-static/ffmpeg /usr/local/bin/ffmpeg && \
    chmod +x /usr/local/bin/ffmpeg && \
    rm -rf /tmp/ffmpeg*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY handler.py .
CMD ["handler.lambda_handler"]
