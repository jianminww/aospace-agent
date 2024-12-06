# Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.20.6-bookworm as builder

WORKDIR /work/

COPY . .

RUN go env -w GO111MODULE=on && make -f Makefile

FROM space-single:local

# space-agent
COPY --from=builder /work/build/aospace /usr/local/bin/aospace

RUN apt-get update \
    && apt-get install -y jq yq

EXPOSE 80 443 5432 6379 3001 2001 8080 5678 5680

HEALTHCHECK --interval=60s --timeout=15s CMD curl -XGET http://localhost:5678/agent/status

# 使用启动脚本作为入口点
ENTRYPOINT ["/usr/local/bin/prestart.sh"]
