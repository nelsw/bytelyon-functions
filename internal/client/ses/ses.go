package ses

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

var (
	source   = "noreply@ByteLyon.com"
	charset  = "UTF-8"
	template = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
        "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office" lang="en">
<head>
    <title></title>
    <meta charset="UTF-8"/>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
    <!--[if !mso]>-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <!--<![endif]-->
    <meta name="x-apple-disable-message-reformatting" content=""/>
    <meta content="target-densitydpi=device-dpi" name="viewport"/>
    <meta content="true" name="HandheldFriendly"/>
    <meta content="width=device-width" name="viewport"/>
    <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no"/>
    <style type="text/css">
        table {
            border-collapse: separate;
            table-layout: fixed;
            mso-table-lspace: 0pt;
            mso-table-rspace: 0pt
        }

        table td {
            border-collapse: collapse
        }

        .ExternalClass {
            width: 100%
        }

        .ExternalClass,
        .ExternalClass p,
        .ExternalClass span,
        .ExternalClass font,
        .ExternalClass td,
        .ExternalClass div {
            line-height: 100%
        }

        body, a, li, p, h1, h2, h3 {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        html {
            -webkit-text-size-adjust: none !important
        }

        body, #innerTable {
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale
        }

        #innerTable img + div {
            display: none;
            display: none !important
        }

        img {
            Margin: 0;
            padding: 0;
            -ms-interpolation-mode: bicubic
        }

        h1, h2, h3, p, a {
            line-height: inherit;
            overflow-wrap: normal;
            white-space: normal;
            word-break: break-word
        }

        a {
            text-decoration: none
        }

        h1, h2, h3, p {
            min-width: 100% !important;
            width: 100% !important;
            max-width: 100% !important;
            display: inline-block !important;
            border: 0;
            padding: 0;
            margin: 0
        }

        a[x-apple-data-detectors] {
            color: inherit !important;
            text-decoration: none !important;
            font-size: inherit !important;
            font-family: inherit !important;
            font-weight: inherit !important;
            line-height: inherit !important
        }

        u + #body a {
            color: inherit;
            text-decoration: none;
            font-size: inherit;
            font-family: inherit;
            font-weight: inherit;
            line-height: inherit;
        }

        a[href^="mailto"],
        a[href^="tel"],
        a[href^="sms"] {
            color: inherit;
            text-decoration: none
        }
    </style>
    <style type="text/css">
        @media (min-width: 481px) {
            .hd {
                display: none !important
            }
        }
    </style>
    <style type="text/css">
        @media (max-width: 480px) {
            .hm {
                display: none !important
            }
        }
    </style>
    <style type="text/css">
        @media (max-width: 480px) {
            .t48 {
                width: 480px !important
            }

            .t42 {
                text-align: center !important
            }

            .t41 {
                vertical-align: top !important;
                width: 600px !important
            }

            .t13, .t18 {
                font-size: 18px !important;
                mso-text-raise: 4px !important
            }

            .t6 {
                font-size: 41px !important;
                mso-text-raise: 3px !important
            }

            .t4 {
                width: 378px !important
            }

            .t32 {
                font-size: 22px !important
            }
        }
    </style>
    <!--[if !mso]>-->
    <link href="https://fonts.googleapis.com/css2?family=Stem:wght@400;500;600;700&display=swap" rel="stylesheet" type="text/css"/>
    <!--<![endif]-->
    <!--[if mso]>
    <xml>
        <o:OfficeDocumentSettings>
            <o:AllowPNG/>
            <o:PixelsPerInch>96</o:PixelsPerInch>
        </o:OfficeDocumentSettings>
    </xml>
    <![endif]-->
</head>
<body id="body" class="t52" style="min-width:100%;Margin:0px;padding:0px;background-color:#FFFFFF;">
<div class="t51" style="background-color:#FFFFFF;">
    <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" align="center">
        <tr>
            <td class="t50" style="font-size:0;line-height:0;mso-line-height-rule:exactly;background-color:#FFFFFF;" valign="top" align="center">
                <!--[if mso]>
                <v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="true" stroke="false">
                    <v:fill color="#FFFFFF"/>
                </v:background>
                <![endif]-->
                <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" align="center" id="innerTable">
                    <tr>
                        <td align="center">
                            <table class="t49" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
                                <tr>
                                    <!--[if mso]>
                                    <td width="600" class="t48" style="background-color:#FFFFFF;width:600px;">
                                    <![endif]-->
                                    <!--[if !mso]>-->
                                    <td class="t48" style="background-color:#FFFFFF;width:600px;">
                                        <!--<![endif]-->
                                        <table class="t47" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;">
                                            <tr>
                                                <td class="t46">
                                                    <div class="t45" style="width:100%;text-align:center;">
                                                        <div class="t44" style="display:inline-block;">
                                                            <table class="t43" role="presentation" cellpadding="0" cellspacing="0" align="center" valign="top">
                                                                <tr class="t42">
                                                                    <td></td>
                                                                    <td class="t41" width="600" valign="top">
                                                                        <table role="presentation" width="100%" cellpadding="0" cellspacing="0" class="t40" style="width:100%;">
                                                                            <tr>
                                                                                <td class="t39" style="background-color:transparent;">
                                                                                    <table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="width:100% !important;">
                                                                                        <tr>
                                                                                            <td>
                                                                                                <div class="t1" style="mso-line-height-rule:exactly;mso-line-height-alt:100px;line-height:60px;font-size:1px;display:block;">
                                                                                                      
                                                                                                </div>
                                                                                            </td>
                                                                                        </tr>
                                                                                        <tr>
                                                                                            <td align="center">
                                                                                                <table class="t5" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
                                                                                                    <tr>
                                                                                                        <!--[if mso]>
                                                                                                        <td width="428"
                                                                                                            class="t4"
                                                                                                            style="width:300px;">
                                                                                                        <![endif]-->
                                                                                                        <!--[if !mso]>-->
                                                                                                        <td class="t4" style="width:300px;">
                                                                                                            <!--<![endif]-->
                                                                                                            <table class="t3" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;">
                                                                                                                <tr>
                                                                                                                    <td class="t2">
                                                                                                                        <div style="font-size:0px;">
                                                                                                                            <img class="t0" style="display:block;border:0;height:auto;width:100%;Margin:0;max-width:100%;" alt="" src="https://bytelyon-public.s3.us-east-1.amazonaws.com/logo.png" mc:edit="img-nH_WVx7XGE1agUolpcAR7u"/>
                                                                                                                        </div>
                                                                                                                    </td>
                                                                                                                </tr>
                                                                                                            </table>
                                                                                                        </td>
                                                                                                    </tr>
                                                                                                </table>
                                                                                            </td>
                                                                                        </tr>
                                                                                        <tr>
                                                                                            <td align="center">
                                                                                                <table class="t11" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
                                                                                                    <tr>
                                                                                                        <!--[if mso]>
                                                                                                        <td width="425"
                                                                                                            class="t10"
                                                                                                            style="width:425px;">
                                                                                                        <![endif]-->
                                                                                                        <!--[if !mso]>-->
                                                                                                        <td class="t10" style="width:425px;">
                                                                                                            <!--<![endif]-->
                                                                                                            <table class="t9" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;">
                                                                                                                <tr>
                                                                                                                    <td class="t8">
                                                                                                                        <div mc:edit="p-npW31agFvVbPxBabOZ2OXv">
                                                                                                                            <h1 class="t6" style="margin:0;Margin:0;font-family:Stem,BlinkMacSystemFont,Segoe UI,Helvetica Neue,Arial,sans-serif;line-height:52px;font-weight:700;font-style:normal;font-size:48px;text-decoration:none;text-transform:none;direction:ltr;color:#0D1D32;text-align:center;mso-line-height-rule:exactly;mso-text-raise:1px;">
                                                                                                                                ByteLyon
                                                                                                                            </h1>
                                                                                                                        </div>
                                                                                                                    </td>
                                                                                                                </tr>
                                                                                                            </table>
                                                                                                        </td>
                                                                                                    </tr>
                                                                                                </table>
                                                                                            </td>
                                                                                        </tr>
                                                                                        <tr>
                                                                                            <td>
                                                                                                <div class="t27" style="mso-line-height-rule:exactly;mso-line-height-alt:40px;line-height:40px;font-size:1px;display:block;">
                                                                                                      
                                                                                                </div>
                                                                                            </td>
                                                                                        </tr>
                                                                                        <tr>
                                                                                            <td align="center">
                                                                                                <table class="t37" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
                                                                                                    <tr>
                                                                                                        <!--[if mso]>
                                                                                                        <td width="308"
                                                                                                            class="t36"
                                                                                                            style="background-color:#2196f3;overflow:hidden;width:308px;border:2px solid #2196f3;border-radius: 4px;">
                                                                                                        <![endif]-->
                                                                                                        <!--[if !mso]>-->
                                                                                                        <td class="t36" style="background-color:#2196f3;overflow:hidden;width:308px;border:2px solid #2196f3;border-radius: 4px;">
                                                                                                            <!--<![endif]-->
                                                                                                            <table class="t35" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;">
                                                                                                                <tr>
                                                                                                                    <td class="t34" style="text-align:center;line-height:58px;mso-line-height-rule:exactly;mso-text-raise:11px;">
                                                                                                                        <div mc:edit="p-nCt8FP5WBQrY_G0y-LQQ8T">
                                                                                                                            <a class="t32" href="{{btn-url}}" style="display:block;margin:0;Margin:0;font-family:Stem,BlinkMacSystemFont,Segoe UI,Helvetica Neue,Arial,sans-serif;line-height:58px;font-weight:600;font-style:normal;font-size:21px;text-decoration:none;direction:ltr;color:#FFFFFF;text-align:center;mso-line-height-rule:exactly;mso-text-raise:11px;" target="_blank">
                                                                                                                                {{btn-txt}}
                                                                                                                            </a>
                                                                                                                        </div>
                                                                                                                    </td>
                                                                                                                </tr>
                                                                                                            </table>
                                                                                                        </td>
                                                                                                    </tr>
                                                                                                </table>
                                                                                            </td>
                                                                                        </tr>
                                                                                        <tr>
                                                                                            <td>
                                                                                                <div class="t38" style="mso-line-height-rule:exactly;mso-line-height-alt:60px;line-height:60px;font-size:1px;display:block;">
                                                                                                      
                                                                                                </div>
                                                                                            </td>
                                                                                        </tr>
                                                                                    </table>
                                                                                </td>
                                                                            </tr>
                                                                        </table>
                                                                    </td>
                                                                    <td></td>
                                                                </tr>
                                                            </table>
                                                        </div>
                                                    </div>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</div>
<div class="gmail-fix" style="display: none; white-space: nowrap; font: 15px courier; line-height: 0;">   
                                   
                           
</div>
</body>
</html>`
)

type Service interface {
	VerifyEmail(ctx context.Context, to, token string) error
	ResetPassword(ctx context.Context, to, token string) error
}

type Client struct {
	*ses.Client
}

func (c *Client) ResetPassword(ctx context.Context, to, token string) error {
	url := fmt.Sprintf("https://ByteLyon.com/reset-password?token=%s", token)
	data := strings.ReplaceAll(template, "{{btn-url}}", url)
	data = strings.ReplaceAll(data, "{{btn-txt}}", "Reset Password")
	subject := "🦁 Reset Password"
	return c.email(ctx, to, subject, data)
}

func (c *Client) VerifyEmail(ctx context.Context, to, token string) error {
	url := fmt.Sprintf("https://ByteLyon.com/verify-email?token=%s", token)
	data := strings.ReplaceAll(template, "{{btn-url}}", url)
	data = strings.ReplaceAll(data, "{{btn-txt}}", "Verify Email")
	subject := "🦁 Verify Email"
	return c.email(ctx, to, subject, data)
}

func (c *Client) email(ctx context.Context, to, subject, data string) error {
	_, err := c.SendEmail(ctx, &ses.SendEmailInput{
		ReplyToAddresses: []string{
			source,
		},
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: &charset,
					Data:    &data,
				},
			},
			Subject: &types.Content{
				Data:    &subject,
				Charset: &charset,
			},
		},
		Source: &source,
	})

	return err
}

func New(ctx context.Context) Service {
	if cfg, err := config.LoadDefaultConfig(ctx); err != nil {
		panic(err)
	} else {
		return &Client{ses.NewFromConfig(cfg)}
	}
}
