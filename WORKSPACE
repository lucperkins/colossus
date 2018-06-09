workspace(name = "colossus")

# Versions
PROMETHEUS_JAVA_VERSION = "0.4.0"

# Imports basic Go rules for Bazel (e.g. go_binary)
git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    commit = "e4d0254fb249a09fb01f052b23d3baddae1b70ec",
)

# Imports the Gazelle tool for Go/Bazel
git_repository(
    name = "bazel_gazelle",
    remote = "https://github.com/bazelbuild/bazel-gazelle",
    commit = "644ec7202aa352b78d65bc66abc2c0616d76cc84",
)

# Imports Docker rules for Bazel (e.g. docker_image)
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.4.0",
)

# Imports gRPC for Java rules (e.g. java_grpc_library)
git_repository(
    name = "io_grpc_grpc_java",
    remote = "https://github.com/grpc/grpc-java",
    tag = "v1.12.0",
)

# Import gRPC for C++
git_repository(
    name = "com_github_grpc_grpc",
    remote = "https://github.com/grpc/grpc.git",
    commit = "17f682d8274ef0b7d1376eeee5e94839a0750e0e",
)

# Import Maven rules for Gradle conversion
git_repository(
    name = "org_pubref_rules_maven",
    remote = "https://github.com/pubref/rules_maven",
    commit = "9c3b07a6d9b195a1192aea3cd78afd1f66c80710",
)

# Loads Maven rules
load("@org_pubref_rules_maven//maven:rules.bzl", "maven_repositories", "maven_repository")

maven_repositories()

# Loads Docker for Java rules (e.g. java_image)
load(
    "@io_bazel_rules_docker//java:image.bzl",
    _java_image_repos = "repositories",
)

_java_image_repos()

# Loads gRPC for Java rules
load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories()

# Loads Go rules for Bazel
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.10.1",
)

# Loads Docker rules for Bazel
load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

# Loads C++ gRPC rules for Bazel
load("@com_github_grpc_grpc//:bazel/grpc_deps.bzl", "grpc_deps")

grpc_deps()

# Loads C++ Docker image rules
load(
    "@io_bazel_rules_docker//cc:image.bzl",
    _cc_image_repos = "repositories",
)

_cc_image_repos()

# Loads Gazelle tool
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# gRPC for Java dependencies (shorthand)
bind(
    name = "grpc-core",
    actual = "@io_grpc_grpc_java//core",
)

bind(
    name = "grpc-netty",
    actual = "@io_grpc_grpc_java//netty",
)

bind(
    name = "grpc-stub",
    actual = "@io_grpc_grpc_java//stub",
)

maven_jar(
    name = "io_prometheus_simpleclient",
    artifact = "io.prometheus:simpleclient:" + PROMETHEUS_JAVA_VERSION,
)

maven_jar(
    name = "io_prometheus_simpleclient_httpserver",
    artifact = "io.prometheus:simpleclient_httpserver:" + PROMETHEUS_JAVA_VERSION,
)

maven_jar(
    name = "me_dinowernli_java_grpc_prometheus",
    artifact = "me.dinowernli:java-grpc-prometheus:0.3.0"
)

# Gazelle-generated Go dependencies
go_repository(
    name = "com_github_inconshreveable_mousetrap",
    commit = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75",
    importpath = "github.com/inconshreveable/mousetrap",
)

go_repository(
    name = "com_github_spf13_cobra",
    commit = "a1f051bc3eba734da4772d60e2d677f47cf93ef4",
    importpath = "github.com/spf13/cobra",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "583c0c0531f06d5278b7d917446061adc344b5cd",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "com_github_go_chi_chi",
    commit = "e83ac2304db3c50cf03d96a2fcd39009d458bc35",
    importpath = "github.com/go-chi/chi",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "c155da19408a8799da419ed3eeb0cb5db0ad5dbc",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "b49d69b5da943f7ef3c9cf91c8777c1f78a0cc3c",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "cbbc999da32df943dac6cd71eb3ee39e1d7838b9",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    commit = "c2828203cd70a50dcccfb2761f8b1f8ceef9a8e9",
    importpath = "github.com/fsnotify/fsnotify",
)

go_repository(
    name = "com_github_hashicorp_hcl",
    commit = "ef8a98b0bbce4a65b5aa4c368430a80ddc533168",
    importpath = "github.com/hashicorp/hcl",
)

go_repository(
    name = "com_github_magiconair_properties",
    commit = "c3beff4c2358b44d0493c7dda585e7db7ff28ae6",
    importpath = "github.com/magiconair/properties",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    commit = "00c29f56e2386353d58c599509e8dc3801b0d716",
    importpath = "github.com/mitchellh/mapstructure",
)

go_repository(
    name = "com_github_pelletier_go_toml",
    commit = "acdc4509485b587f5e675510c4f2c63e90ff68a8",
    importpath = "github.com/pelletier/go-toml",
)

go_repository(
    name = "com_github_spf13_afero",
    commit = "63644898a8da0bc22138abf860edaf5277b6102e",
    importpath = "github.com/spf13/afero",
)

go_repository(
    name = "com_github_spf13_cast",
    commit = "8965335b8c7107321228e3e3702cab9832751bac",
    importpath = "github.com/spf13/cast",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    commit = "7c0cea34c8ece3fbeb2b27ab9b59511d360fb394",
    importpath = "github.com/spf13/jwalterweatherman",
)

go_repository(
    name = "com_github_spf13_viper",
    commit = "b5e8006cbee93ec955a89ab31e0e3ce3204f3736",
    importpath = "github.com/spf13/viper",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "5420a8b6744d3b0345ab293f6fcba19c978f1183",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "org_golang_x_text",
    commit = "f21a4dfb5e38f5895301dc265a8def02365cc3d0",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "com_github_go_pg_pg",
    commit = "5b73ce88484575f3480edf393237f6bf79d5f166",
    importpath = "github.com/go-pg/pg",
)

go_repository(
    name = "com_github_jinzhu_inflection",
    commit = "04140366298a54a039076d798123ffa108fff46c",
    importpath = "github.com/jinzhu/inflection",
)

go_repository(
    name = "com_github_golang_protobuf",
    commit = "b4deda0973fb4c70b50d226b1af49f3da59f5265",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "org_golang_google_genproto",
    commit = "86e600f69ee4704c6efbf6a2a40a5c10700e76c2",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "d11072e7ca9811b1100b80ca0269ac831f06d024",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "org_golang_x_net",
    commit = "5f9ae10d9af5b1c89ae6904293b14b064d4ada23",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "com_github_go_redis_redis",
    commit = "0f9028adf0837cf93c9705817493e5f6997cf026",
    importpath = "github.com/go-redis/redis",
)

go_repository(
    name = "com_github_unrolled_render",
    commit = "65450fb6b2d3595beca39f969c411db8f8d5c806",
    importpath = "github.com/unrolled/render",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "3a771d992973f24aa725d07868b467d1ddfceafb",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "c5b7fccd204277076155f10851dad72b76a49317",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "99fa1f4be8e564e8a6b613da7fa6f46c9edafc6c",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "7600349dcfe1abd18d72d3a1770870d9800a7801",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "94663424ae5ae9856b40a9f170762b4197024661",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_prometheus",
    commit = "c225b8c3b01faf2899099b768856a9e916e5087b",
    importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
)
