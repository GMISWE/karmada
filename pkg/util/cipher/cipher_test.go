/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/15
 * @Desc    : 加密解密测试
 */

package cipher

import (
	"testing"
)

var (
	kb1 = ``
	kb2 = ``
	kb3 = `#!/bin/bash


init() {
  mkdir -p {{.MOUNT_POINT}} || true

  CURL_PATH=$(which curl 2>/dev/null || echo "")
  if [ -z "$CURL_PATH" ]; then
    apt update
    apt install -y curl
  else
    echo "curl already installed at $CURL_PATH"
  fi

  FUSE_PATH=$(which fusermount 2>/dev/null || echo "")
  if [ -z "$FUSE_PATH" ]; then
    apt update
    apt install -y fuse
  else
    echo "fuse already installed at $FUSE_PATH"
  fi

  # check gcloud cli
  GCLOUD_PATH=$(which gcloud 2>/dev/null || echo "")
  if [ -z "$GCLOUD_PATH" ]; then
    echo "gcloud cli is not installed, installing it..."
    # add google cloud sdk source
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

    # install apt-transport-https ca-certificates gnupg
    apt-get install -y apt-transport-https ca-certificates gnupg
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

    # update and install cloud sdk
    apt-get update
    apt-get install -y google-cloud-sdk
  else
    echo "gcloud cli already installed at $GCLOUD_PATH"
  fi

  # create gcloud config dir
  if [ ! -d ~/.config/gcloud ]; then
    mkdir -p ~/.config/gcloud
  fi

  if [ ! -f ~/.config/gcloud/application_default_credentials.json ]; then
    # create application_default_credentials.json
    cat <<EOF > ~/.config/gcloud/application_default_credentials.json
{{.GCP_APPLICATION_DEFAULT_CREDENTIALS}}
EOF
  fi
}

mount() {
  init
  # install JuiceFS client
  if [ ! -f "{{.JUICEFS_PATH}}" ]; then
    echo "JuiceFS client not found, downloading..."
    curl -L {{.JUICEFS_CONSOLE_HOST}}/onpremise/juicefs -o {{.JUICEFS_PATH}} && chmod +x {{.JUICEFS_PATH}}
  else
    # get current installed JuiceFS version
    CURRENT_VERSION=$({{.JUICEFS_PATH}} version 2>/dev/null | awk '{print $3}')
    # compare version
    if [ "$CURRENT_VERSION" != "{{.JUICEFS_VERSION}}" ]; then
      echo "JuiceFS version mismatch (current: $CURRENT_VERSION, required: {{.JUICEFS_VERSION}}), downloading new version..."
      curl -L {{.JUICEFS_CONSOLE_HOST}}/onpremise/juicefs -o {{.JUICEFS_PATH}} && chmod +x {{.JUICEFS_PATH}}
      {{.JUICEFS_PATH}} version -u
    fi
  fi

  {{.JUICEFS_PATH}} mount --cache-dir {{.JUICEFS_CACHE_DIR}} --token {{.JUICEFS_TOKEN}} {{.JUICEFS_NAME}} {{.MOUNT_POINT}} {{.MOUNT_OPTIONS}}

}

umount() {
  {{.JUICEFS_PATH}} umount {{.MOUNT_POINT}}
  rm -rf {{.MOUNT_POINT}}
  # safe remove log file
  if [ -f "{{.STORAGE_PATH}}" ]; then
    rm -f "{{.STORAGE_PATH}}"
    echo "storage path removed: {{.STORAGE_PATH}}"
  fi
  echo "{{.JUICEFS_NAME}} umount done"
}

dry_run() {
  echo "dry run"
  init
  mount
  umount
}

main(){
  case $1 in
    mount) mount ;;
    umount) umount ;;
    dry-run) dry_run ;;
    *) echo "Usage: $0 {mount|umount|dry-run}" ;;
  esac
}

main "$@"`
)

func TestEncrypt(t *testing.T) {
	cases := []struct {
		Name string
		text []byte
	}{
		{"a", []byte(kb1)},
		{"b", []byte(kb2)},
		{"c", []byte(kb3)},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans, err := EncryptCompact(c.text, "uOvKLmVfztaXGpNYd4Z0I1SiT7MweJhl"); err != nil {
				t.Fatalf("encrypt text %s failed: %+v",
					c.text, err)
			} else {
				t.Logf("%s encrypt text is { %s }", c.Name, ans)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	cases := []struct {
		Name string
		text string
	}{
		{"a", `KAl-BcH6vzo1WQ4WogEQPcKvEnRgolYNJRNwuudvModdTpaSj338T4QtdElnHUPiY9X8tT5VG__7EJv9x9VO8D75SjrEAecFaWtWEdAEwW0297Cl6quYS47dUVfCcDO0Xsmt8j6im5Ac8u4rfYIo-22hRmc0TEOcwemsQh26PMUxsHoiDkedAETUwzNLVWzMDPE8vJX0WRycyM-vAm3euKo9Zq2Ijk73mLSrFvWrSpsUYyxKadsOPQM7XlM491qLO95xGwupiaIdje2th7MhrasoxO9jE31QV4o0gzfKCyWRDlny4CSyKVaQStLK9YrML7BXNcadlUuAM97eaGxx6m7IWAxvCGXltmoTWPTIpYFNwxgEyZqtYyUWzB4KShNX3TKL7mu2VPUw9giyVO2Ys9xPj5KfCOlYtP4DdkH1FnVOY-Mb94X3khlqF2fVOl-S7cQwgsRsa_QsMz1gd4iu1EQ-4CIqslEaU_AmiZi1IPWjJQRg3rL_OXcMZ5RGLry3`},
		{"b", "HlZO9LEf_QCcd8gMTKTjNndMz9Cw7d51lDke14sbrWZvNM-LfpRjZpakfv53xLve9q_XN4LcMmST88NdCkxM6vO5lUTt-QiWnsy8NfiVWpCgLsVQ7s4fcL8I2bgKo9VyD0P075lOg3Z2qJGGzvr6Yp5jYN1Be1QDvmhn_W2L80FKUZ-1aO9s2hFkb0fLE-fKmhiOaAjDp-kmrMB0XfKhfavGM1-LWnkC4E7XdouFa8E3qli5V0L8bsPBb0cLf6YiFQQY7Wj_pkykyt4uceHnb7kV31W4nemO7Rz5BNEQgyVFQjw2CDAQHYZE0JBo3IbPfgBwrQuNAaYF0RVRS2P-2v2PxcZ-CnFZJcvYbULjn1QtbXTg1RkjYLXoDDV3vYaQ_0Rp7-YNHKy6Kn1rIEHE3SNxqhH_P7NXQLL-BYavv15oCMoDXSqsrY_tjTo3JwIwW8hocMVzLBQ5Hh0w-ZoI1-4PTjdsMh3KRWoCT0EIqcmKtzRUDTdyzZkYzvym6nIPyOfjurtD-85pbiBzgVHFE6lUE-jDdAd6Q5wzyUxROUNNSxlW45yAVL-OiQwgYkN-YQnDlTesjvkZBgq8ZqFW5r9q1jSmNLr54ngmoS37H5GAL529g_-7C3tass7N3zWNJ8UjPBuMbFvSRZ6GZZx0CESjh_Kkmo_ABsPHg7_HDgSir9TbqCVhGoVGN3UFseXE4MTT4lTKFd2EHsiJ-5f669pxzW57M0vpwGWgRL_Cdv24h2jtn6NrYyJPo-fqmEOt3Cl4dFxbw1GAxeoElSckLCw8gfpTHjVlVusQiX6bYBVT3WJKPrBpxycR7xONnhTd5yzjTlMfONXQbiH9yvTqfC9oM_P31hr3ooPgVAB42krc5sIsNjZ8jZzpuZH5Tl_FbmAhzapqNcwu-YhFD6ZhoU5yHKf3aXG4DmV1-4efq6NxEdvd0Fd5IyKFtbLX2Fev7vdhUcu_CRMY-IkmKTH_tJ6aETeVF_2g_-W5XqJci7ncnS2Ia8n1JPWQk2ZMoP6RKBx_rRKlRaChMbU6Sl6VYgmWCBI6FH1p-Tg4jfcaOjYx2Tg-E4kEti_bT_gPt_uFycHVCvOxPaxHE7dmPbfEu4tyFVCHjXEPucYDf0_bka3KAfD9GpFE1MPH-hw3SgFyLYnXt71AxNd5BoLGBrUg0RoKsGrxdYUbEiu_GnEEiBB3cfIGsXnX2mMEThkbeHQrbWdVaJ4YVUr5AmN6cVFDa6ZTh-xNPkp4r8zjc-CQgYg8rdNPfOm-B3dwemwErqcGihBKAexcFne0Fa7BB7HIEpIoHzOdb9ZiC8pIjcg4MeErGjL4B3N1RxdeiIi_qg-wFpMq-YvwRBW2p2S3zS49DDXbxU-DeH8qGjHL7qmx_TL0bzvC7iaFGcxpUOBE7SmHLbVhFzdLtEnIoy6-pvushO8Y0Y7lHnPHpdm3C7oeAyiyYsHqn_OcMohfUAj8SNZS_n3zKwopheCIRyGBEdnmcZGnmnXhl13CuhnsJ3H9iHlVeNHs44Y8LMQDU_JZxORoGVyP8wPxfMmw2F6dMjeXx_lJkcqIQtkle4Oc42nHsMOeiNmLKf6IC8Em3qsOtmbq2lPsDOXwPtTmJ3_nuU7E60qyomlolJIE1FhPOYUBl_9adCl1Enlt2ACpV2eamqRUTJxIDjljzJoHGH4BeuVjVg00fcNx0zU7tPmjVnatNoj0tjxMYQ-MISrh_F81XYQgtRKlF-s1QvQguk1B6mGwteGwyY0AjUAyJw54puw-i6JBOyBDT03aab2xTW10vx0Qma0gX3YZgOQejJ77AgVmJIHEvANHCFxmUXEfs69xR6D2i-iyNNrpwNFEdeHkDNQ7mZ49SVd0HFgtqnZHxKpXZee13Q_C3wfThx66EGFpAGKz9hrvK7bGw7VLwhqYbQulHdN3BIQAjqEMZTlZ0c_lVIDZnj7Kf-meTNjPWWF8eFAoYv_Zo0L96LktUOSZud8VUpE3-fYqRNxpQk0FEixweo4-K-NjwvKCpkLQrqps9mWPMoxLN9tH8RKmXJlBqp37C_0-lg0B6oTP_QXMqt-VOS0qtUwV3AOvLCQKaP8_TLKL9O0WWjGjJb_fM0k-qZu9DrLsICb0EwMRbz1nDDJaZUi-nzh5uDrAZaI5Wka3SmvQcAoOFi4catKc-RvArjbcbyIjFQDC635aC7v4ULfwQYpvGIYduMd_L95KpM0vEIqIuqZwf7igygmBVAlz4c67WWLtWdK8gTcQtlwEAHq_Vbs8BkL3VpS9KS8sJTKXRl1qNDD1uakTyqbqmqcZu_gwQBj_G6zh6tJtQX_PoY817T2M3KqC5Rid3e2tLYWR8UuK6zeYSbYGfP7niGPA0A6LcRPvG1P7I9v1NjKSbnLF1yQIOQ2KQYy8RSHPfPeYxCgOzobSy-odp9zmgRGdP5cH7oG6mMJC2NE4SHP5CT87c0to2l-3P5bkF_X_YDLDi7cBpVtGD2_yJZ6gJdsx7FdK7xIJbRx3rMMhATjHqPa5ZL-x8_dgXQiSMxheBdIyDu0IuH-NJ0QWmz-qFfQ93JV9Jd1YjllmQcNSHV6n_159FITzDlriMEsVBy_M1D_BPRlbsfMr87Dk4_IlOqALucm-SgcR3cF_XUa7TlZlUlDHrGg6g6j-BV8KXVOav3wtEZq1qAiRZm1_tQDoEaw68tys9Zqmg3A_6Rn_15PwGmEVn3Fj6P1GiKB-cNxxzYm5DbCcWINujwnLrOtJnd5yhpxAo5JUtm39eLQaEMKCjlx7iph89OxGngWTk7kP-eeQqLAe0SpOf9B2WD2neBNucrBCXAShIACDllBmNgAYBhEOrSWWn3w4idJYLlXGBHTRk8mxeHQDt5sKzb90obvJSgF9ML4r5iYWolXxNSYQ9j_kqlKhbviZJA_vQo5GrNOPHoycpCwaEF7ci-ff_QRdAqLyaPt9DftE1DtyueTWiZKuxoxR9hnt7cIh64WJovt-RykB4scsZ-FE8ZJ649wLsEks7UtroxlS6SlYr2mI9e5Lzf6P9lTkjz_uiO9Q04O7sjJv09mDhYnHA1k_R32YeyYWkBlfDP4KpC1AXo_HndRyIE1xv7ji8pE1Bv9i3GXmLnUG1dE-nOaphXiASMtZqQxFsnj_Ti3R_3_9LdJ8tK4nEqzLekKlnLtw1wkDAwCzJq4-E64AfDIqGj0uKVgjycNIflvi0IJrazEV55rUCKfmTvbKBM_SwJMjHFuNB9aimN9KhyO5Rbif_j93quBQ6csROR4BhLzYEWw9nyP-4fCpOi-tvgMtYrphb9C35QmSS9KoWfVclaPL_FVdNhl15azmZBdG2HFe6IfBHKWnmEygRK2Tr3D39TZmfVCFvDePsDKl9NPMO4h0vrJRIl-pOm6bPhvE1YJyhUn4YyHK-HBaRNPnwqO6XCE8fgnD1Bm1pBIsZwF8jh9gxLL9moXWTAMHFXBdN1-dC0XlFKYIIFyNIr2dqPiNFQRkR3IpZJvGaPhI-2GFZrjPpj-B3c5oC6bi"},
		{"c", "ho1s5qTr+11w9t/9W/d6yg=="},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans, err := DecryptCompact(c.text, "uOvKLmVfztaXGpNYd4Z0I1SiT7MweJhl"); err != nil {
				t.Fatalf("decrypt text %s failed: %+v",
					c.text, err)
			} else {
				t.Logf("decrypt is %s", ans)
			}
		})
	}
}
